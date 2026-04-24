package main

import (
	"log"

	"github.com/meow-Kil/KeyboardsGO/internal/adapters/storage"
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server"
	"github.com/meow-Kil/KeyboardsGO/internal/core/service"
)

func main() {
	_storage, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer _storage.Close()

	_keyboardService := service.NewKeyboard(_storage)
	_keycapTypeService := service.NewKeycapType(_storage)   
	srv := server.New(_keyboardService, _keycapTypeService, _storage, "./static")
	srv.Listen()
}
