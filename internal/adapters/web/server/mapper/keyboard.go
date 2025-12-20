package mapper

import ("github.com/meow-Kil/KeyboardsGO/internal/core/domain"
		"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server/dto"
)
		

func ToDto(keyboard *domain.Keyboard) *dto.Keyboard {
	return &dto.Keyboard{
		Id : keyboard.Id,
		KeycapType : keyboard.KeycapType,
		BaseType : keyboard.BaseType,
		SwitchType : keyboard.SwitchType,
		Color : keyboard.Color,
	}
}
func FromDto(keyboard *dto.Keyboard) *domain.Keyboard {
	return &domain.Keyboard{
		Id : keyboard.Id,
		KeycapType : keyboard.KeycapType,
		BaseType : keyboard.BaseType,
		SwitchType : keyboard.SwitchType,
		Color : keyboard.Color,
	}
}

func ToDtoList(keyboard []domain.Keyboard) []dto.Keyboard {
	var slice = make([]dto.Keyboard, len(keyboard))
		for i, a:= range keyboard{
			slice[i]= *ToDto(&a)
		}
	return slice
	}