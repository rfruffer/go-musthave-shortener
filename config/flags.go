package config

import (
	"flag"
	"os"
)

type Config struct {
	StartHost  string
	ResultHost string
}

func ParseFlags() *Config {
	startHost := flag.String("a", "localhost:8080", "address and port to run server")
	resultHost := flag.String("b", "http://localhost:8080", "base URL for shortened links")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		*startHost = envRunAddr
	}
	if envResultHost := os.Getenv("BASE_URL"); envResultHost != "" {
		*resultHost = envResultHost
	}

	return &Config{
		StartHost:  *startHost,
		ResultHost: *resultHost,
	}
}
