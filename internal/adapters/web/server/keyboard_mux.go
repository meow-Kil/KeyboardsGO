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

func (s *Server) MuxKeycapType() {
	
	s.mux.HandleFunc("GET /keycap_types", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Требуется авторизация"})
			return
		}
		types := s.keycapType.GetAll()
		s.writeJson(w, mapper.ToKeycapTypeDtoList(types))
	})


	s.mux.HandleFunc("GET /keycap_types/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Требуется авторизация"})
			return
		}
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		kt := s.keycapType.Get(id)
		if kt == nil {
			w.WriteHeader(http.StatusNotFound)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Тип не найден"})
			return
		}
		s.writeJson(w, mapper.ToKeycapTypeDto(kt))
	})


	s.mux.HandleFunc("POST /keycap_types", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Требуются права администратора"})
			return
		}
		var req dto.KeycapType
		if err := s.parseObject(w, r, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Некорректные данные"})
			return
		}
		created := s.keycapType.Create(*mapper.FromKeycapTypeDto(&req))
		s.writeJson(w, map[string]interface{}{
			"success": true,
			"message": "Тип кейкапа добавлен",
			"data":    mapper.ToKeycapTypeDto(created),
		})
	})


	s.mux.HandleFunc("PUT /keycap_types/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Требуются права администратора"})
			return
		}
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		var req dto.KeycapType
		if err := s.parseObject(w, r, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Некорректные данные"})
			return
		}
		updated := s.keycapType.Update(id, *mapper.FromKeycapTypeDto(&req))
		if updated == nil {
			w.WriteHeader(http.StatusNotFound)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Тип не найден"})
			return
		}
		s.writeJson(w, map[string]interface{}{
			"success": true,
			"message": "Тип обновлён",
			"data":    mapper.ToKeycapTypeDto(updated),
		})
	})


	s.mux.HandleFunc("DELETE /keycap_types/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			s.writeJson(w, map[string]interface{}{"success": false, "error": "Требуются права администратора"})
			return
		}
		id := s.parseId(w, r)
		if id == 0 {
			return
		}
		s.keycapType.Delete(id)
		s.writeJson(w, map[string]interface{}{"success": true, "message": "Тип удалён"})
	})
}