package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/crypto"
	"github.com/CyrilSbrodov/passManager.git/server/internal/handlers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/storage/repositories"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/client/postgres"
)

type ServerApp struct {
	router *chi.Mux
	cfg    *config.Config
	logger *loggers.Logger
	crypto crypto.RSA
}

func NewServerApp() *ServerApp {
	cfg := config.ConfigInit()
	logger := loggers.NewLogger()
	router := chi.NewRouter()
	c := crypto.NewRSA(*cfg)

	return &ServerApp{
		router: router,
		cfg:    cfg,
		logger: logger,
		crypto: *c,
	}
}

func (a *ServerApp) Run() {
	client, err := postgres.NewClient(context.Background(), 5, a.cfg, a.logger)
	if err != nil {
		a.logger.LogErr(err, "")
		os.Exit(1)
	}
	store, err := repositories.NewStore(client, a.cfg, a.logger)
	if err != nil {
		a.logger.LogErr(err, "")
		os.Exit(1)
	}
	var tokenAuth *jwtauth.JWTAuth
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	handler := handlers.NewHandler(store, a.logger, a.crypto, tokenAuth)

	//регистрация хендлера
	handler.Register(a.router)

	srv := http.Server{
		Addr:    a.cfg.Addr,
		Handler: a.router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.LogErr(err, "server not started")
		}
	}()
	a.logger.LogInfo("server is listen:", a.cfg.Addr, "start server")

	//gracefullshutdown
	<-done

	a.logger.LogInfo("", "", "server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.LogErr(err, "Server Shutdown Failed")
	}
	a.logger.LogInfo("", "", "Server Exited Properly")
}
