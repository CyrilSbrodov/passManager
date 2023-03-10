package config

type Config struct {
	Addr             string
	CryptoPROKey     string
	CryptoPROKeyPath string
}

func ConfigInit() *Config {
	return &Config{}
}
