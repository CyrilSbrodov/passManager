// Package auth пакет предоставляет возможность создания и проверки токенов авторизации.
package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt      = "fsdfsdfsdfsdf"
	signinKey = "gsdfgdfg$564gdf"
	tokenTTL  = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

func GenerateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Login: id,
	})
	return token.SignedString([]byte(signinKey))
}

func ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signinKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not type *tokenClaims")
	}
	return claims.Login, nil
}
