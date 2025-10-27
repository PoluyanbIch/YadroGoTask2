package config

import (
	"os"
)

type Config struct {
	GRPCPORT string
}

func Load() *Config {
	port := ":" + os.Getenv("WORDS_GRPC_PORT")
	if port == "" {
		port = ":8080"
	}
	return &Config{GRPCPORT: port}
}
