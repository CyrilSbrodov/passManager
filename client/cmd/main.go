package main

import "github.com/CyrilSbrodov/passManager.git/client/app"

// Main - функция сборки и запуска клиента.
func main() {
	client := app.NewApp()
	client.Run()
}
