# gRPC поддержка для сервиса сокращения URL

## Описание

В сервис добавлена поддержка gRPC протокола наряду с существующим HTTP API. Все хендлеры доступны через gRPC и функционируют идентично HTTP версии.

## Архитектура

### Общий слой бизнес-логики

Создан `CommonURLService` который объединяет всю бизнес-логику и используется как HTTP, так и gRPC хендлерами:

- `internal/services/commonService.go` - общий сервис для обоих протоколов
- `internal/grpc/server.go` - gRPC сервер и хендлеры  
- `internal/grpc/types.go` - типы данных для gRPC
- `internal/grpc/service.go` - интерфейсы gRPC сервиса

### Конфигурация

Добавлены новые параметры конфигурации:

- Флаг `-g` или переменная окружения `GRPC_ADDRESS` для указания адреса gRPC сервера (по умолчанию `0.0.0.0:3200`)
- Поддержка в JSON конфигурации через поле `grpc_address`

## Доступные методы gRPC

1. **CreateShortURL** - создание короткой ссылки
2. **CreateShortJSONURL** - создание короткой ссылки (JSON формат)
3. **GetShortURL** - получение оригинального URL
4. **Ping** - проверка подключения к БД
5. **Batch** - пакетное создание коротких ссылок
6. **GetUserURLs** - получение всех URL пользователя
7. **BatchDelete** - пакетное удаление URL
8. **GetStats** - получение статистики сервиса

## Запуск

### Локальный запуск

```bash
go run ./cmd/shortener -g=0.0.0.0:3200
```

### Docker

```bash
# Сборка образа
docker build -t shortener-grpc -f Dockerfile.grpc .

# Запуск с docker-compose
docker-compose -f docker-compose.grpc.yml up
```

## Порты

- **8080** - HTTP сервер
- **3200** - gRPC сервер (по умолчанию)

## Переменные окружения

- `SERVER_ADDRESS` - адрес HTTP сервера
- `GRPC_ADDRESS` - адрес gRPC сервера  
- `BASE_URL` - базовый URL для коротких ссылок
- `DATABASE_DSN` - строка подключения к PostgreSQL
- `SECRET_KEY` - секретный ключ для JWT
- `TRUSTED_SUBNET` - доверенная подсеть для админ методов

## Примеры использования

### Тестирование HTTP API (остается без изменений)

```bash
# Создание короткой ссылки
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Получение статистики
curl http://localhost:8080/api/internal/stats
```

### Общий слой бизнес-логики

HTTP и gRPC хендлеры используют один и тот же `CommonURLService`, что обеспечивает:

- Единообразное поведение между протоколами
- Переиспользование кода
- Легкость добавления новых протоколов
- Централизованную обработку ошибок

## Преимущества архитектуры

1. **Фасад-паттерн**: HTTP и gRPC хендлеры являются фасадами к общему коду бизнес-логики
2. **DRY принцип**: Нет дублирования логики между протоколами
3. **Легкая поддержка**: Изменения в бизнес-логике автоматически применяются к обоим протоколам
4. **Тестируемость**: Общий слой можно тестировать независимо от протоколов
5. **Расширяемость**: Легко добавить поддержку новых протоколов (например, GraphQL)

## Файлы проекта

```
├── api/
│   └── shortener.proto          # Протобуф определения
├── cmd/shortener/
│   └── main.go                  # Запуск HTTP и gRPC серверов
├── config/
│   └── flags.go                 # Конфигурация с поддержкой gRPC
├── internal/
│   ├── grpc/
│   │   ├── server.go           # gRPC сервер
│   │   ├── types.go            # Типы данных
│   │   └── service.go          # Интерфейсы
│   └── services/
│       └── commonService.go    # Общий слой бизнес-логики
├── tests/
│   └── common_service_test.go  # Тесты общего сервиса
├── Dockerfile.grpc             # Docker образ с gRPC
└── docker-compose.grpc.yml     # Docker Compose конфигурация
```

## Тестирование

```bash
# Тестирование общего сервиса
go test ./tests -run TestCommonService -v

# Компиляция проекта
go build ./cmd/shortener
```