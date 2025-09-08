// Package config парсит конфигурацию из флагов и переменных окружения.
package config

import (
	"flag"
	"os"
)

// Config предоставляет переменные для конфигурации сервисов.
type Config struct {
	StartHost   string
	ResultHost  string
	FilePath    string
	DBDSN       string
	Storage     string
	SecretKey   string
	EnableHTTPS bool
	CertFile    string
	KeyFile     string
}

// ParseFlags устанавливает пути
func ParseFlags() *Config {
	startHost := flag.String("a", "0.0.0.0:8080", "address and port to run server")
	resultHost := flag.String("b", "http://localhost:8080", "base URL for shortened links")
	filePath := flag.String("f", "", "path to file storage")
	dbDSN := flag.String("d", "", "database DSN for PostgreSQL")
	enableHTTPS := flag.Bool("s", false, "enable HTTPS") // новый флаг
	certFile := flag.String("cert", "./certs/cert.pem", "path to TLS certificatefile")
	keyFile := flag.String("key", "./certs/key.pem", "path to TLS private keyfile")
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "verysecretkey"
	}

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
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		*dbDSN = envDB
	}
	if os.Getenv("ENABLE_HTTPS") == "true" {
		*enableHTTPS = true
	}
	if envCertFile := os.Getenv("CERT_FILE"); envCertFile != "" {
		*certFile = envCertFile
	}
	if envKeyFile := os.Getenv("KEY_FILE"); envKeyFile != "" {
		*keyFile = envKeyFile
	}

	storage := ""
	if *dbDSN != "" {
		storage = "postgres"
	}

	return &Config{
		StartHost:   *startHost,
		ResultHost:  *resultHost,
		FilePath:    *filePath,
		DBDSN:       *dbDSN,
		Storage:     storage,
		SecretKey:   secretKey,
		EnableHTTPS: *enableHTTPS,
		CertFile:    *certFile,
		KeyFile:     *keyFile,
	}
}
