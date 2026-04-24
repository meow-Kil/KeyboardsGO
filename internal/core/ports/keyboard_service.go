package ports

import "github.com/meow-Kil/KeyboardsGO/internal/core/domain"

type KeyboardService interface {
	GetAll() []domain.Keyboard
	Get(id uint) *domain.Keyboard
	New(keyboard domain.Keyboard) *domain.Keyboard
	Delete(id uint)
	Update(id uint, keyboard domain.Keyboard)  *domain.Keyboard
}

type KeycapTypeService interface {
	GetAll() []domain.KeycapType
	Get(id uint) *domain.KeycapType
	Create(kt domain.KeycapType) *domain.KeycapType
	Update(id uint, kt domain.KeycapType) *domain.KeycapType
	Delete(id uint)
}