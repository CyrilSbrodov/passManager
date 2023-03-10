package main

import "github.com/CyrilSbrodov/passManager.git/client/app"

func main() {
	client := app.NewApp()
	client.Run()
}
