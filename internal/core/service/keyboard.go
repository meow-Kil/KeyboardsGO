package service

import (
	"github.com/meow-Kil/KeyboardsGO/internal/core/domain"
	"github.com/meow-Kil/KeyboardsGO/internal/core/ports"
)

type Keyboard struct{
	storage ports.Storage
}
func NewKeyboard(storage ports.Storage) *Keyboard {
	return &Keyboard{storage}
}
func(a *Keyboard) GetAll() []domain.Keyboard {
	return a.storage.Get()
}
func(a *Keyboard) Get(id uint) *domain.Keyboard {
	return a.storage.GetById(id)
}

func(a *Keyboard) New(keyboard domain.Keyboard) *domain.Keyboard {
	keyboard = a.storage.Add(keyboard)
	return &keyboard
}


func(a *Keyboard) Delete(id uint)  {
	a.storage.Remove(id)
}

func(a *Keyboard) Update(id uint, keyboard domain.Keyboard)  *domain.Keyboard {
	return a.storage.Update(id,keyboard)
}