package main

import "github.com/CyrilSbrodov/passManager.git/server/internal/app"

func main() {
	srv := app.NewServerApp()
	srv.Run()
}
