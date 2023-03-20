package main

import "github.com/CyrilSbrodov/passManager.git/server/internal/app"

// main функция сборки и запуска сервера.
func main() {
	srv := app.NewServerApp()
	srv.Run()
}
