package ports

import "github.com/meow-Kil/KeyboardsGO/internal/core/domain"

type Storage interface {
	
	Add(a domain.Keyboard) domain.Keyboard
	Get() []domain.Keyboard
	GetById(id uint) *domain.Keyboard
	Remove(id uint)
	Update(id uint, keyboard domain.Keyboard) *domain.Keyboard

	
	AddUser(login, password string, isAdmin bool) (*domain.User, error)
	GetUserByLogin(login string) (*domain.User, error)

	
	AddKeycapType(kt domain.KeycapType) domain.KeycapType
	GetKeycapTypes() []domain.KeycapType
	GetKeycapTypeByID(id uint) *domain.KeycapType
	UpdateKeycapType(id uint, kt domain.KeycapType) *domain.KeycapType
	DeleteKeycapType(id uint)
}