package storage
   

import "github.com/meow-Kil/KeyboardsGO/internal/core/domain"

 
type Storage struct{
	keyboards []domain.Keyboard
	nextId uint
}

var stor Storage

func New() *Storage {
	stor= Storage{keyboards: make([]domain.Keyboard, 0, 1024)}
	return &stor
}
func (s *Storage) Add(a domain.Keyboard) domain.Keyboard {
	s.nextId++
	a.Id = s.nextId
	s.keyboards = append(s.keyboards, a)
	return a
}
func (s *Storage) Get() []domain.Keyboard {
	return s.keyboards
}
func (s *Storage) GetById(id uint) *domain.Keyboard {
	for _, a := range s.keyboards{
		if a.Id == id {
			b:=a 
			return &b
		}
	}
	return nil
}
func (s *Storage) Remove(id uint) {
	for i, a:= range s.keyboards{
		if a.Id == id{
			s.keyboards = append(s.keyboards[:i] ,s.keyboards[i+1:]...)
			return
		}
	}
}
func (s *Storage) Update(id uint, keyboard domain.Keyboard)  *domain.Keyboard {
	for i, a := range s.keyboards{
		if a.Id == id{
			keyboard.Id = id
			s.keyboards[i] = keyboard
		}
	}
	return &keyboard
}