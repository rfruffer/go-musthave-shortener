// Package config парсит конфигурацию из флагов и переменных окружения.
package config

import (
	"encoding/json"
	"flag"
	"os"
)

// JSONConfig структура для JSON конфигурации
type JSONConfig struct {
	ServerAddress   string `json:"server_address,omitempty"`
	BaseURL         string `json:"base_url,omitempty"`
	FileStoragePath string `json:"file_storage_path,omitempty"`
	DatabaseDSN     string `json:"database_dsn,omitempty"`
	EnableHTTPS     *bool  `json:"enable_https,omitempty"`
}

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

// loadJSONConfig загружает конфигурацию из JSON файла
func loadJSONConfig(configPath string) (*JSONConfig, error) {
	if configPath == "" {
		return &JSONConfig{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config JSONConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// ParseFlags устанавливает пути
func ParseFlags() *Config {
	startHost := flag.String("a", "0.0.0.0:8080", "address and port to run server")
	resultHost := flag.String("b", "http://localhost:8080", "base URL for shortened links")
	filePath := flag.String("f", "", "path to file storage")
	dbDSN := flag.String("d", "", "database DSN for PostgreSQL")
	enableHTTPS := flag.Bool("s", false, "enable HTTPS")
	certFile := flag.String("cert", "./certs/cert.pem", "path to TLS certificatefile")
	keyFile := flag.String("key", "./certs/key.pem", "path to TLS private keyfile")
	configPath := flag.String("c", "", "path to JSON config file")
	_ = flag.String("config", "", "path to JSON config file (alias for -c)")
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "verysecretkey"
	}

	flag.Parse()

	// Обработка алиаса -config
	if flag.Lookup("config").Value.String() != "" && *configPath == "" {
		*configPath = flag.Lookup("config").Value.String()
	}

	// Проверяем переменную окружения CONFIG
	if envConfigPath := os.Getenv("CONFIG"); envConfigPath != "" && *configPath == "" {
		*configPath = envConfigPath
	}

	// Загружаем JSON конфигурацию (наименьший приоритет)
	jsonConfig, err := loadJSONConfig(*configPath)
	if err != nil {
		// Если файл указан, но не найден - это ошибка
		if *configPath != "" {
			panic("Failed to load config file: " + err.Error())
		}
	}

	// Устанавливаем значения по умолчанию из JSON (наименьший приоритет)
	if jsonConfig.ServerAddress != "" && *startHost == "0.0.0.0:8080" {
		*startHost = jsonConfig.ServerAddress
	}
	if jsonConfig.BaseURL != "" && *resultHost == "http://localhost:8080" {
		*resultHost = jsonConfig.BaseURL
	}
	if jsonConfig.FileStoragePath != "" && *filePath == "" {
		*filePath = jsonConfig.FileStoragePath
	}
	if jsonConfig.DatabaseDSN != "" && *dbDSN == "" {
		*dbDSN = jsonConfig.DatabaseDSN
	}
	if jsonConfig.EnableHTTPS != nil && !*enableHTTPS {
		*enableHTTPS = *jsonConfig.EnableHTTPS
	}

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
