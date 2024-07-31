package natsbroker

import "os"

type Config struct {
	URL           string
	StreamName    string
	FilterSubject string
}

func NewDefaultConfig() *Config {
	return &Config{
		URL:           "nats://localhost:4223",
		StreamName:    "L0_STEAM",
		FilterSubject: "l0.*",
	}
}

func NewEnvConfig() *Config {
	config := &Config{}
	if url, isExists := os.LookupEnv("NATS_URL"); isExists {
		config.URL = url
	} else {
		config.URL = "nats://localhost:4223"
	}

	if stream, isExists := os.LookupEnv("NATS_STREAMNAME"); isExists {
		config.StreamName = stream
	} else {
		config.StreamName = "L0_STREAM"
	}

	if subject, isExists := os.LookupEnv("NATS_STREAM_SUBJECT"); isExists {
		config.FilterSubject = subject
	} else {
		config.FilterSubject = "l0.*"
	}

	return config
}
