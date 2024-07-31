package apiserver

import "os"

type Config struct {
	BindAddr string
}

func NewDefaultConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}

func NewEnvConfig() *Config {
	config := &Config{}
	if bindAddr, isExists := os.LookupEnv("APISERVER_BINDADDR"); isExists {
		config.BindAddr = bindAddr
	} else {
		config.BindAddr = ":8080"
	}

	return config
}
