package server

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/meow-Kil/KeyboardsGO/internal/core/ports"
	"strconv"
)

type Server struct{
	srv *http.Server
	mux *http.ServeMux
	keyboard ports.KeyboardService
}

func parseObject[T any](w http.ResponseWriter, r *http.Request) *T {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil 
	}
	return &obj
}


func (s *Server) parseId (w http.ResponseWriter, r *http.Request) uint {
	id := r.PathValue("id")
	_id, err := strconv.ParseUint(id,10,64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 0
	}
	return uint(_id)
}
func (s *Server) writeJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(jsonData)
}



func New(keyboard ports.KeyboardService) *Server {
	mux := http.NewServeMux()
	

	srv := http.Server{
		Addr: ":9090",
		Handler: mux, 
		DisableGeneralOptionsHandler: false,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		keyboard: keyboard,
		srv: &srv,
		mux: mux,
	}
}


func (s *Server) Listen() {
	s.MuxKeyboard()
	err := s.srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}