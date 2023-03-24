// Package handlers позволяет получать данные от клиентов, обрабатывать и отправлять в репозиторий для дальнейшей обработки.
// Данный модуль дает возможность обрабатывать бинарные данные.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
)

func (h *Handler) CollectBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)

		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var c models.CryptoBinaryData

		if err := json.Unmarshal(content, &c); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.CollectBinary(&c, userID)

		switch statusCode {
		case http.StatusOK:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			return
		case http.StatusInternalServerError:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}
}

func (h *Handler) GetBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string)

		statusCode, data, err := h.Storage.GetBinary(userID)

		switch statusCode {
		case http.StatusNoContent:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNoContent)
			rw.Write([]byte(err.Error()))
			return
		case http.StatusInternalServerError:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		dJSON, err := json.Marshal(data)
		if err != nil {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(dJSON)
	}
}

func (h *Handler) DeleteBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)

		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var data models.CryptoBinaryData

		if err := json.Unmarshal(content, &data); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.DeleteBinary(&data, userID)

		switch statusCode {
		case http.StatusOK:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			return
		case http.StatusInternalServerError:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}
}

func (h *Handler) UpdateBinary() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)

		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var c models.CryptoBinaryData

		if err := json.Unmarshal(content, &c); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.UpdateBinary(&c, userID)

		switch statusCode {
		case http.StatusOK:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			return
		case http.StatusInternalServerError:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}
}
