package postgresql

import "os"

type Config struct {
	URL string
}

func NewDefaultConfig() *Config {
	return &Config{
		URL: "postgresql://localdbuser:localdbuserpass@localhost:5433/wbl0",
	}
}

func NewEnvConfig() *Config {
	config := &Config{}
	if url, isExists := os.LookupEnv("POSTGRESQL_URL"); isExists {
		config.URL = url
	} else {
		config.URL = "postgresql://localdbuser:localdbuserpass@localhost:5433/wbl0"
	}

	return config
}
