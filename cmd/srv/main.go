package main

import (
	"fmt"
	"os"

	"github.com/meow-Kil/KeyboardsGO/internal/adapters/storage"
	"github.com/meow-Kil/KeyboardsGO/internal/adapters/web/server"
	"github.com/meow-Kil/KeyboardsGO/internal/core/service"
)

func main() {
	_storage, err := storage.NewPostgresStorage()
	if err != nil {}
	defer _storage.Close()

	_keyboardService := service.NewKeyboard(_storage)
	pwd, _ := os.Getwd()
fmt.Println(pwd)
srv := server.New(_keyboardService, "./static")
	srv.Listen()
}