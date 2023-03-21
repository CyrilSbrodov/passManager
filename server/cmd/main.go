package main

import (
	"fmt"

	"github.com/CyrilSbrodov/passManager.git/server/internal/app"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main функция сборки и запуска сервера.
func main() {
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	srv := app.NewServerApp()
	srv.Run()
}
