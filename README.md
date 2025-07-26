# USDT Rate Service

Микросервис для получения и хранения курсов валют USDT через gRPC API. Сервис получает данные глубины рынка от биржи Grinex и предоставляет актуальные курсы покупки и продажи.

## Возможности

- ⚡ gRPC API для получения курсов валют
- 🏛️ Интеграция с биржей Grinex
- 🗄️ Хранение данных в PostgreSQL
- 🔄 Валидация данных глубины рынка
- 📊 Структурированное логирование с Zap
- 🐳 Docker контейнеризация
- 🧪 Юнит-тесты с моками

## Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │    │  Rates Service  │    │   Grinex API    │
│                 │───▶│                 │───▶│                 │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                               ▼
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │                 │
                    └─────────────────┘
```

### Компоненты

- **gRPC Handler** - обработка входящих запросов
- **Rates Service** - бизнес-логика получения и сохранения курсов
- **Depth Provider** - адаптер для получения данных с Grinex
- **Repository** - слой доступа к данным PostgreSQL

## Требования

- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose
- Protocol Buffers compiler (protoc)

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd usdt-rate-service
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Генерация Protocol Buffers

```bash
make proto
```

### 4. Запуск через Docker Compose

```bash
make run
# или
docker-compose -f build/docker-compose.yml up --build
```

Сервис будет доступен на порту `50052`.

### 5. Локальная разработка

```bash
# Запуск только PostgreSQL
docker-compose -f build/docker-compose.yml up postgres migrations

# Установка переменных окружения
export GRPC_ADDRESS=:50051
export DATABASE_ADDRESS="postgres://usdt_user:usdt_password@localhost:5432/usdt_rates?sslmode=disable"
export GRINEX_ADDRESS="https://grinex.io"
export LOG_LEVEL=debug

# Сборка и запуск
make build
./bin/usdt-rate-service
```

## API

### gRPC методы

#### GetRates
Получение актуального курса для указанного рынка.

```protobuf
rpc GetRates(GetRatesRequest) returns (GetRatesResponse);

message GetRatesRequest {
  string market = 1;  // например, "usdtrub"
}

message GetRatesResponse {
  Rate rate = 1;
}

message Rate {
  string askPrice = 1;   // цена продажи
  string bidPrice = 2;   // цена покупки  
  int64 timestamp = 3;   // временная метка
}
```

#### HealthCheck
Проверка состояния сервиса.

```protobuf
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
```

### Пример использования

```bash
# Используя grpcurl
grpcurl -plaintext -d '{"market":"usdtrub"}' localhost:50052 rates.RatesService/GetRates

# Результат
{
  "rate": {
    "askPrice": "102.50",
    "bidPrice": "102.30", 
    "timestamp": 1737901234
  }
}
```

## Конфигурация

Конфигурация через переменные окружения или флаги командной строки:

| Переменная окружения | Флаг | Описание | Пример |
|---------------------|------|----------|---------|
| `GRPC_ADDRESS` | `-grpc-addr` | Адрес gRPC сервера | `:50051` |
| `DATABASE_ADDRESS` | `-db-addr` | Строка подключения к БД | `postgres://user:pass@host:port/db` |
| `GRINEX_ADDRESS` | `-grinex-addr` | URL API Grinex | `https://grinex.io` |
| `LOG_LEVEL` | `-log-level` | Уровень логирования | `debug`, `info`, `warn`, `error` |

## Разработка

### Структура проекта

```
usdt-rate-service/
├── api/proto/          # Protocol Buffers определения
├── build/              # Docker Compose и миграции
├── cmd/                # Точка входа приложения
├── config/             # Конфигурация
├── internal/           # Внутренний код
│   ├── adapter/        # Адаптеры внешних сервисов
│   ├── handler/grpc/   # gRPC обработчики
│   ├── infra/grinex/   # Клиент для Grinex API
│   ├── models/         # Модели данных
│   ├── pb/             # Сгенерированные Protocol Buffers
│   ├── repository/     # Слой доступа к данным
│   ├── server/grpc/    # gRPC сервер
│   └── service/        # Бизнес-логика
├── mocks/              # Моки для тестирования
└── pkg/                # Общие утилиты
```

### Команды Make

```bash
make proto        # Генерация Protocol Buffers
make mocks        # Генерация моков
make build        # Сборка приложения
make test         # Запуск тестов
make lint         # Линтинг кода
make docker-build # Сборка Docker образа
make run          # Запуск через Docker Compose
```

### Генерация моков

```bash
# Установка mockery
go install github.com/vektra/mockery/v2@latest

# Генерация моков
make mocks
```

### Тестирование

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
go test -cover ./...

# Запуск конкретного теста
go test -run TestRatesService_GetRates ./internal/service/
```

## База данных

### Миграции

Миграции находятся в `build/migrations/postgres/` и выполняются автоматически при запуске через Docker Compose.

### Схема таблицы rates

```sql
CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    market VARCHAR(20) NOT NULL,
    ask DECIMAL(20,8) NOT NULL,
    bid DECIMAL(20,8) NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rates_market_timestamp ON rates(market, timestamp);
```

## Мониторинг и логирование

Сервис использует структурированное логирование с библиотекой Zap. Логи включают:

- Информацию о запросах и ответах
- Ошибки интеграции с внешними API
- Метрики производительности
- Отладочную информацию (при уровне debug)

Пример лога:
```json
{
  "level": "info",
  "ts": 1737901234.123,
  "caller": "service/rates.go:65",
  "msg": "Rate saved successfully",
  "service": "RatesService",
  "method": "GetRates",
  "market": "usdtrub",
  "rate": {
    "market": "usdtrub",
    "ask": "102.50",
    "bid": "102.30",
    "timestamp": 1737901234
  }
}
```

## Производство

### Docker образ

```bash
# Сборка образа
make docker-build

# Запуск контейнера
docker run -p 50051:50051 \
  -e GRPC_ADDRESS=:50051 \
  -e DATABASE_ADDRESS="postgres://..." \
  -e GRINEX_ADDRESS="https://grinex.io" \
  -e LOG_LEVEL=info \
  usdt-rate-service:latest
```

### Рекомендации для продакшена

- Использовать connection pooling для PostgreSQL
- Настроить healthcheck endpoints
- Добавить метрики Prometheus
- Настроить graceful shutdown
- Использовать TLS для gRPC соединений
- Добавить rate limiting
- Настроить мониторинг и alerting

## Лицензия

[MIT License](LICENSE)

## Контакты

Для вопросов и предложений создавайте issues в репозитории.