package config

import (
	"flag"
)

type Config struct {
	StartHost  string
	ResultHost string
}

func ParseFlags() *Config {
	startHost := flag.String("a", "localhost:8080", "address and port to run server")
	resultHost := flag.String("b", "http://localhost:8080", "base URL for shortened links")

	flag.Parse()

	return &Config{
		StartHost:  *startHost,
		ResultHost: *resultHost,
	}
}
