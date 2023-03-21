// Package config - это пакет с конфигом, который позволяет гибко запускать приложение.
package config

import "flag"

// Config - структура конфига.
type Config struct {
	Addr             string `json:"address" env:"ADDRESS"`
	CryptoPROKey     string `json:"crypto_key" env:"CRYPTO_KEY"`
	CryptoPROKeyPath string `json:"crypto_key_path" env:"CRYPTO_KEY_PATH"`
}

// ConfigInit - инициализация конфига.
func ConfigInit() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.CryptoPROKey, "crypto-key", "private.pem", "path to file")
	flag.StringVar(&cfg.CryptoPROKeyPath, "crypto-key-path", "./client/crypto/", "path to folder")
	return cfg
}
