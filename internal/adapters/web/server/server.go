package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"github.com/meow-Kil/KeyboardsGO/internal/core/ports"
)

type Server struct {
	srv      *http.Server
	mux      *http.ServeMux
	keyboard ports.KeyboardService
	storage  ports.Storage
	sessions map[string]Session
	keycapType ports.KeycapTypeService
}

type Session struct {
	UserID   uint
	Login    string
	IsAdmin  bool
	Expires  time.Time
}

func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) writeJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) parseId(w http.ResponseWriter, r *http.Request) uint {
	id := r.PathValue("id")
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 0
	}
	return uint(_id)
}

func (s *Server) parseObject(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func (s *Server) isAdminAuthorized(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}
	
	token := parts[1]
	session, exists := s.sessions[token]
	if !exists {
		return false
	}
	
	if time.Now().After(session.Expires) {
		delete(s.sessions, token)
		return false
	}
	
	return session.IsAdmin
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	type RegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}

	var req RegisterRequest
	if err := s.parseObject(w, r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error": "Некорректный запрос",
		})
		return
	}


	existingUser, err := s.storage.GetUserByLogin(req.Login)
	if err == nil && existingUser != nil {
		w.WriteHeader(http.StatusConflict)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error": "Пользователь с таким логином уже существует",
		})
		return
	}

	user, err := s.storage.AddUser(req.Login, req.Password, req.IsAdmin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error": "Ошибка регистрации: " + err.Error(),
		})
		return
	}

	s.writeJson(w, map[string]interface{}{
		"success": true,
		"message": "Пользователь успешно зарегистрирован",
		"user_id": user.ID,
		"login": user.Login,
		"is_admin": user.IsAdmin,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		LoginType string `json:"login_type,omitempty"`
	}

	var req LoginRequest
	if err := s.parseObject(w, r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error":   "Некорректный запрос",
		})
		return
	}

	user, err := s.storage.GetUserByLogin(req.Login)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error":   "Пользователь не найден",
		})
		return
	}
	
	if user == nil || user.Password != req.Password {
		w.WriteHeader(http.StatusUnauthorized)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error":   "Неверный пароль",
		})
		return
	}

	if req.LoginType == "admin" && !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		s.writeJson(w, map[string]interface{}{
			"success": false,
			"error":   "У вас нет прав доступа к админ-панели",
		})
		return
	}

	
	token := s.generateSessionToken()
	

	s.sessions[token] = Session{
		UserID:  user.ID,
		Login:   user.Login,
		IsAdmin: user.IsAdmin,
		Expires: time.Now().Add(24 * time.Hour),
	}

	responseData := map[string]interface{}{
		"success":  true,
		"user_id":  user.ID,
		"login":    user.Login,
		"is_admin": user.IsAdmin,
		"token":    token,
		"message":  "Авторизация успешна",
	}

	if req.LoginType == "admin" {
		responseData["redirect"] = "/admin.html"
	} else {
		responseData["redirect"] = "/user.html"
	}

	s.writeJson(w, responseData)
}

func (s *Server) generateSessionToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatInt(time.Now().Unix(), 36)
}

func New(keyboard ports.KeyboardService, keycapType ports.KeycapTypeService, storage ports.Storage, staticPath string) *Server {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticPath))
	mux.Handle("/", fileServer)

	srv := http.Server{
		Addr:    ":1000",
		Handler: mux,
	}

	server := &Server{
		keyboard: keyboard,
		keycapType: keycapType,
		storage:  storage,
		srv:      &srv,
		mux:      mux,
		sessions: make(map[string]Session),
	}

	return server
}


func (s *Server) Listen() {
	s.MuxKeyboard()
	s.MuxKeycapType()
	s.srv.Handler = s.enableCORS(s.mux)
	err := s.srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) isAuthorized(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}
	token := parts[1]
	session, exists := s.sessions[token]
	if !exists || time.Now().After(session.Expires) {
		if exists {
			delete(s.sessions, token)
		}
		return false
	}
	return true
}