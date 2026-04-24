package mapper

import (
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/dto"
	"github.com/meow-Kil/KeyboardsGO/internal/core/domain"
)

func ToKeycapTypeDto(kt *domain.KeycapType) *dto.KeycapType {
	if kt == nil {
		return nil
	}
	return &dto.KeycapType{
		ID:   kt.ID,
		Name: kt.Name,
	}
}

func FromKeycapTypeDto(kt *dto.KeycapType) *domain.KeycapType {
	return &domain.KeycapType{
		ID:   kt.ID,
		Name: kt.Name,
	}
}

func ToKeycapTypeDtoList(types []domain.KeycapType) []dto.KeycapType {
	list := make([]dto.KeycapType, len(types))
	for i, v := range types {
		list[i] = *ToKeycapTypeDto(&v)
	}
	return list
}