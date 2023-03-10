package app

import (
	"net/http"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/client/crypto"
	"github.com/CyrilSbrodov/passManager.git/client/manager"
)

type App struct {
	crypto  crypto.Crypto
	manager manager.Managers
	client  http.Client
	logger  loggers.Logger
	cfg     config.Config
}

func NewApp() *App {
	c := crypto.NewRSA()
	m := manager.NewManager()
	client := &http.Client{}
	return &App{
		crypto:  c,
		manager: m,
		client:  *client,
	}
}

func (a *App) Run() {

}
