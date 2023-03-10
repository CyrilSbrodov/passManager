package handlers

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/storage"
)

type Handlers interface {
	Register(router *chi.Mux)
}

type Handler struct {
	storage.Storage
	logger       loggers.Logger
	sessionStore sessions.Store
}

func NewHandler(storage storage.Storage, logger *loggers.Logger, sessionStore sessions.Store) Handlers {
	return &Handler{
		storage,
		*logger,
		sessionStore,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	compressor := middleware.NewCompressor(gzip.DefaultCompression)
	r.Use(compressor.Handler)
	r.Post("/api/user/register", h.Registration())
	r.Post("/api/user/login", h.Login())

	r.Group(func(r chi.Router) {
		r.Use(h.Auth)
		r.Post("/api/user/data", h.CollectData())
		r.Get("/api/user/data", h.GetAll())
		r.Post("/api/user/data/update", h.UpdateData())
		r.Post("/api/user/data/delete", h.DeleteData())
	})
}

func (h *Handler) Registration() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	})
}

func (h *Handler) CollectData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (h *Handler) UpdateData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (h *Handler) DeleteData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}
