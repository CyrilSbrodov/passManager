package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/CyrilSbrodov/passManager.git/server/pkg/auth"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			fmt.Println("empty auth header")
			http.Error(rw, "empty auth header", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			fmt.Println("invalid auth header")
			http.Error(rw, "invalid auth header", http.StatusUnauthorized)
			return
		}
		userID, err := auth.ParseToken(headerParts[1])
		if err != nil {
			fmt.Println("invalid parse token")
			http.Error(rw, "invalid parse token", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(ctx, "user_id", userID))

		next.ServeHTTP(rw, r)
	})
}
