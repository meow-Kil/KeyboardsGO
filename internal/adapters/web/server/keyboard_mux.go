package server

import (
	"net/http" 
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/mapper"
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/dto"
)

func (s Server) MuxKeyboard() {
	s.mux.HandleFunc("GET /keyboard", func(w http.ResponseWriter, r *http.Request){
		s.writeJson(w, mapper.ToDtoList(s.keyboard.GetAll()))
	})
	s.mux.HandleFunc("GET /keyboard/{id}", func(w http.ResponseWriter, r *http.Request){
		id := s.parseId(w,r)
		if id > 0 {
			s.writeJson(w,mapper.ToDto(s.keyboard.Get(id)))
		}
		
	})
	s.mux.HandleFunc("POST /keyboard", func(w http.ResponseWriter, r *http.Request){
		keyboard := parseObject[dto.Keyboard](w,r)
		if keyboard == nil {
			return
		}
		s.writeJson(w, mapper.ToDto(s.keyboard.New(*mapper.FromDto(keyboard))))
	})
	s.mux.HandleFunc("PUT /keyboard/{id}", func(w http.ResponseWriter, r *http.Request){
		id:= s.parseId(w,r)
		if id == 0 {
			return
		}
		keyboard := parseObject[dto.Keyboard](w,r)
		if keyboard == nil {
			return
		}
		s.writeJson(w, s.keyboard.Update(id, *mapper.FromDto(keyboard)))
	})
	s.mux.HandleFunc("DELETE /keyboard/{id}", func(w http.ResponseWriter, r *http.Request){
		id := s.parseId(w,r)
		if id > 0 {
			s.keyboard.Delete(id)
			s.writeJson(w, struct{}{})
		}
		
	})
}