package server

import (
	"net/http"

	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/dto"
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/mapper"
)

func (s *Server) MuxKeyboard() {
	s.mux.HandleFunc("POST /register", s.handleRegister)
	s.mux.HandleFunc("POST /login", s.handleLogin)

	s.mux.HandleFunc("GET /keyboard", func(w http.ResponseWriter, r *http.Request) {
		s.writeJson(w, mapper.ToDtoList(s.keyboard.GetAll()))
	})

	s.mux.HandleFunc("GET /keyboard/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		keyboard := s.keyboard.Get(id)
		if keyboard == nil {
			w.WriteHeader(http.StatusNotFound)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error": "Клавиатура не найдена",
			})
			return
		}
		s.writeJson(w, mapper.ToDto(keyboard))
	})

	s.mux.HandleFunc("POST /keyboard", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error":   "Требуется авторизация администратора",
			})
			return
		}
		
		var keyboard dto.Keyboard
		if err := s.parseObject(w, r, &keyboard); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error": "Некорректные данные",
			})
			return
		}
		
		newKeyboard := s.keyboard.New(*mapper.FromDto(&keyboard))
		s.writeJson(w, map[string]interface{}{
			"success": true,
			"message": "Клавиатура успешно добавлена",
			"data": mapper.ToDto(newKeyboard),
		})
	})

	s.mux.HandleFunc("PUT /keyboard/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error":   "Требуется авторизация администратора",
			})
			return
		}
		
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		
		var keyboard dto.Keyboard
		if err := s.parseObject(w, r, &keyboard); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error": "Некорректные данные",
			})
			return
		}
		
		updatedKeyboard := s.keyboard.Update(id, *mapper.FromDto(&keyboard))
		s.writeJson(w, map[string]interface{}{
			"success": true,
			"message": "Клавиатура успешно обновлена",
			"data": mapper.ToDto(updatedKeyboard),
		})
	})

	s.mux.HandleFunc("DELETE /keyboard/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{
				"success": false,
				"error":   "Требуется авторизация администратора",
			})
			return
		}
		
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		
		s.keyboard.Delete(id)
		s.writeJson(w, map[string]interface{}{
			"success": true,
			"message": "Клавиатура успешно удалена",
		})
	})
}