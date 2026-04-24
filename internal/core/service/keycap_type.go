package service

import (
	"github.com/meow-Kil/KeyboardsGO/internal/core/domain"
	"github.com/meow-Kil/KeyboardsGO/internal/core/ports"
)

type KeycapType struct {
	storage ports.Storage
}

func NewKeycapType(storage ports.Storage) *KeycapType {
	return &KeycapType{storage: storage}
}

func (s *KeycapType) GetAll() []domain.KeycapType {
	return s.storage.GetKeycapTypes()
}

func (s *KeycapType) Get(id uint) *domain.KeycapType {
	return s.storage.GetKeycapTypeByID(id)
}

func (s *KeycapType) Create(kt domain.KeycapType) *domain.KeycapType {
	created := s.storage.AddKeycapType(kt)
	return &created
}

func (s *KeycapType) Update(id uint, kt domain.KeycapType) *domain.KeycapType {
	return s.storage.UpdateKeycapType(id, kt)
}

func (s *KeycapType) Delete(id uint) {
	s.storage.DeleteKeycapType(id)
}