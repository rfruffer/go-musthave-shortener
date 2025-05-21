package config

import (
	"flag"
	"os"
)

type Config struct {
	StartHost  string
	ResultHost string
	FilePath   string
}

func ParseFlags() *Config {
	startHost := flag.String("a", "localhost:8080", "address and port to run server")
	resultHost := flag.String("b", "http://localhost:8080", "base URL for shortened links")
	filePath := flag.String("f", "storage.json", "path to file storage")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		*startHost = envRunAddr
	}
	if envResultHost := os.Getenv("BASE_URL"); envResultHost != "" {
		*resultHost = envResultHost
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		*filePath = envFilePath
	}

	return &Config{
		StartHost:  *startHost,
		ResultHost: *resultHost,
		FilePath:   *filePath,
	}
}
