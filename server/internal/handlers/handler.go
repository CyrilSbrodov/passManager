package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/crypto"
	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
	"github.com/CyrilSbrodov/passManager.git/server/internal/storage"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/auth"
)

type Handlers interface {
	Register(router *chi.Mux)
}

type Handler struct {
	storage.Storage
	crypto crypto.RSA
	logger loggers.Logger
	token  *jwtauth.JWTAuth
}

func NewHandler(storage storage.Storage, logger *loggers.Logger, c crypto.RSA, token *jwtauth.JWTAuth) Handlers {
	return &Handler{
		Storage: storage,
		crypto:  c,
		logger:  *logger,
		token:   token,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	compressor := middleware.NewCompressor(gzip.DefaultCompression)
	r.Use(compressor.Handler)
	r.Group(func(r chi.Router) {
		r.Post("/api/register", h.Registration())
		r.Post("/api/login", h.Login())
	})
	r.Group(func(r chi.Router) {
		r.Use(h.userIdentity)
		r.Post("/api/data/cards", h.CollectCards())
		r.Post("/api/data/text", h.CollectText())
		r.Post("/api/data/password", h.CollectPassword())
		r.Post("/api/data/binary", h.CollectBinary())

		r.Get("/api/data/cards", h.GetCards())
		r.Get("/api/data/text", h.GetText())
		r.Get("/api/data/password", h.GetPasswords())
		r.Get("/api/data/binary", h.GetBinary())

		r.Post("/api/data/delete/cards", h.DeleteCards())
		r.Post("/api/data/delete/text", h.DeleteText())
		r.Post("/api/data/delete/password", h.DeletePassword())
		r.Post("/api/data/delete/binary", h.DeleteBinary())

		r.Post("/api/data/update/cards", h.UpdateCards())
		r.Post("/api/data/update/text", h.UpdateText())
		r.Post("/api/data/update/password", h.UpdatePassword())
		r.Post("/api/data/update/binary", h.UpdateBinary())
	})
}

func (h *Handler) Registration() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u models.User
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		if err := json.Unmarshal(content, &u); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		if u.Login == "" || u.Password == "" {
			rw.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("login or password is empty")
			rw.Write([]byte(err.Error()))
			return
		}
		id, err := h.Storage.Register(&u)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusConflict)
			err = fmt.Errorf("login %v is already registered", u.Login)
			rw.Write([]byte(err.Error()))
			return
		}
		token, err := auth.GenerateToken(id)
		if err != nil {
			h.logger.LogErr(err, "error")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(err.Error()))
			return
		}
		var sender models.KeyAndToken
		sender.Key = h.crypto.Public
		sender.Token = token

		send, err := json.Marshal(sender)
		if err != nil {
			fmt.Println(err)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(send)
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u models.User
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(content, &u); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		if u.Login == "" || u.Password == "" {
			rw.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("login or password is empty")
			rw.Write([]byte(err.Error()))
			return
		}

		id, err := h.Storage.Login(&u)
		if err != nil {
			h.logger.LogErr(err, "wrong password or login")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(err.Error()))
			return
		}
		token, err := auth.GenerateToken(id)
		if err != nil {
			h.logger.LogErr(err, "error")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(err.Error()))
			return
		}

		var sender models.KeyAndToken
		sender.Key = h.crypto.Public
		sender.Token = token

		send, err := json.Marshal(sender)
		if err != nil {
			fmt.Println(err)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(send)
	}
}
