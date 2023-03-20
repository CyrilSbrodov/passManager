package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
)

func (h *Handler) CollectCards() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var c models.CryptoCard

		if err := json.Unmarshal(content, &c); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.CollectCard(&c, userID)

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

func (h *Handler) GetCards() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string)

		statusCode, data, err := h.Storage.GetCards(userID)

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

func (h *Handler) DeleteCards() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)

		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var data models.CryptoCard

		if err := json.Unmarshal(content, &data); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.DeleteCard(&data, userID)

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

func (h *Handler) UpdateCards() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)

		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var c models.CryptoCard

		if err := json.Unmarshal(content, &c); err != nil {
			h.logger.LogErr(err, "failed to unmarshal data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		userID := r.Context().Value("user_id").(string)

		statusCode, err := h.Storage.UpdateCard(&c, userID)

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
