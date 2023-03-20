// Package config пакет конфига сервера.
package config

import "flag"

// Config - структура конфига.
type Config struct {
	Addr             string `json:"address" env:"ADDRESS"`
	DatabaseDSN      string `json:"database_dsn" env:"DATABASE_DSN"`
	CryptoPROKey     string `json:"crypto_key" env:"CRYPTO_KEY"`
	CryptoPROKeyPath string `json:"crypto_key_path" env:"CRYPTO_KEY_PATH"`
	SessionKey       string `env:"SESSION_KEY"`
}

// ConfigInit - инициализация конфига.
func ConfigInit() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "DATABASE_DSN")
	flag.StringVar(&cfg.SessionKey, "k", "secret", "session key")
	flag.StringVar(&cfg.CryptoPROKey, "crypto-key", "private.pem", "path to file")
	flag.StringVar(&cfg.CryptoPROKeyPath, "crypto-key-path", "", "path to folder")
	return cfg
}
