package config

type Config struct {
	Addr             string
	DataBaseURI      string
	CryptoPROKey     string
	CryptoPROKeyPath string
}

func ConfigInit() *Config {
	return &Config{}
}
