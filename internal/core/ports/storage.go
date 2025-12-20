package ports

import "github.com/meow-Kil/KeyboardsGO/internal/core/domain"

type Storage interface{
	Add(a domain.Keyboard) domain.Keyboard
	Get() []domain.Keyboard
	GetById(id uint) *domain.Keyboard
	Remove(id uint)
	Update(id uint, keyboard domain.Keyboard)  *domain.Keyboard
}