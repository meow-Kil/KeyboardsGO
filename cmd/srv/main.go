package main

import (
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/storage"
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server"
	"github.com/meow-Kil/KeyboardsGO/internal/core/service"
)

func main() {
	_storage := storage.New()
	_keyboardService := service.NewKeyboard(_storage)
	srv := server.New(_keyboardService)
	srv.Listen()
}
