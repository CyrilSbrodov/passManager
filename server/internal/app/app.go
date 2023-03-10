package app

import (
	"crypto/rsa"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/crypto"
)

type ServerApp struct {
	router  *chi.Mux
	cfg     config.Config
	logger  *loggers.Logger
	Crypto  crypto.Crypto
	private *rsa.PrivateKey
}

func NewServerApp() *ServerApp {
	return &ServerApp{}
}

func (a *ServerApp) Run() {
}
