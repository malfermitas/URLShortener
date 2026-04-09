# URL Shortener

Высокопроизводительный сервис для сокращения ссылок, написанный на Go 1.25 с использованием гексагональной архитектуры. Проект включает в себя полноценный REST API, современный Web UI, систему кэширования, расширенный мониторинг и трассировку.

## 🚀 Особенности

- **Гексагональная архитектура (Hexagonal Architecture / Ports & Adapters)**: Четкое разделение бизнес-логики, интерфейсов и внешних зависимостей.
- **Высокая производительность**: Использование Redis для кэширования сокращенных ссылок.
- **Web UI**: Простой и удобный интерфейс для создания ссылок и просмотра аналитики.
- **Аналитика**: Отслеживание количества переходов и метаданных запросов.
- **Наблюдаемость (Observability)**:
  - **Метрики**: Prometheus и Grafana (готовые дашборды).
  - **Трассировка**: OpenTelemetry и Jaeger для отслеживания пути запросов.
  - **Логирование**: Структурированные логи.
- **Устойчивость**: Реализация паттерна Retry для работы с базой данных.
- **Нагрузочное тестирование**: Настроенные сценарии для Vegeta.

## 🛠 Технологический стек

- **Язык**: Go 1.25
- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **База данных**: PostgreSQL
- **Кэш**: Redis
- **Мониторинг**: Prometheus, Grafana
- **Трассировка**: Jaeger (OpenTelemetry)
- **Нагрузочное тестирование**: Vegeta
- **Контейнеризация**: Docker, Docker Compose

## 📁 Структура проекта

```text
URLShortener/
├── cmd/
│   └── app/                 # Точка входа в приложение
├── internal/
│   ├── adapter/             # Адаптеры (внешние слои)
│   │   ├── in/              # Входящие (REST API, Web UI)
│   │   └── out/             # Исходящие (Postgres, Redis, Generator)
│   ├── core/                # Бизнес-логика
│   │   ├── model/           # Доменные модели
│   │   ├── port/            # Порты (интерфейсы)
│   │   └── service/         # Реализация сервисов
│   ├── config/              # Конфигурация приложения
│   ├── logging/             # Настройка логгера
│   ├── metrics/             # Экспорт метрик (Prometheus)
│   ├── tracing/             # Настройка трассировки (Jaeger)
│   └── migrations/          # SQL миграции для БД
├── grafana/                 # Конфигурация и дашборды Grafana
├── prometheus/              # Конфигурация Prometheus
├── loadtest/                # Инструменты для нагрузочного тестирования
├── Dockerfile               # Docker-образ приложения
└── docker-compose.yml       # Оркестрация всех сервисов
```

## 🚦 Быстрый запуск

### С использованием Docker Compose (рекомендуется)

Самый простой способ запустить весь стек (App, Postgres, Redis, Prometheus, Grafana, Jaeger):

```bash
docker-compose up -d --build
```

После запуска:
- **Web UI**: [http://localhost:8080](http://localhost:8080)
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Grafana**: [http://localhost:3000](http://localhost:3000) (логин/пароль: `admin/admin`)
- **Jaeger UI**: [http://localhost:16686](http://localhost:16686)

### Локальный запуск (для разработки)

1. Убедитесь, что у вас установлены Go 1.25, PostgreSQL и Redis.
2. Настройте `config.yaml` или переменные окружения (см. `internal/config/config.go`).
3. Запустите приложение:
   ```bash
   go run cmd/app/main.go
   ```

## 📡 API Endpoints

### 1. Сокращение ссылки
- **URL**: `POST /shorten`
- **Body** (JSON):
  ```json
  {
    "origin_url": "https://example.com/very/long/url",
    "custom_url": "my-cool-link" (опционально)
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "short_url": "http://localhost:8080/s/my-cool-link"
  }
  ```

### 2. Переход по ссылке
- **URL**: `GET /s/:short_url`
- Выполняет 302 Redirect на оригинальный URL и записывает событие в аналитику.

### 3. Аналитика ссылки
- **URL**: `GET /s/:short_url/analytics`
- Возвращает статистику переходов по конкретной ссылке.

## 📊 Мониторинг и диагностика

### Grafana
В проект включен готовый дашборд для Grafana (`grafana/dashboards/urls_dashboard.json`), который автоматически импортируется при запуске через Docker Compose. Он отображает:
- RPS (Requests Per Second)
- HTTP Error Rate
- Статистику сокращенных ссылок

### Jaeger (Трассировка)
Все запросы проходят через OpenTelemetry middleware. Вы можете увидеть детализацию выполнения каждого запроса (время работы с БД, кэшем и т.д.) в интерфейсе Jaeger.

## ⚡ Нагрузочное тестирование

Для проверки производительности в папке `loadtest` есть скрипты:

```bash
cd loadtest
./run_loadtest.sh
```

Скрипт использует утилиту **Vegeta** для генерации нагрузки и создает HTML-отчет в папке `loadtest/results`.

## ⚙️ Конфигурация

Параметры приложения можно задать через файл `config.yaml` или переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_PORT` | Порт сервера | `8080` |
| `DATABASE_HOST` | Хост PostgreSQL | `localhost` |
| `DATABASE_NAME` | Имя БД | (обязательно) |
| `REDIS_ADDRESS` | Адрес Redis | `localhost:6379` |
| `LOG_LEVEL` | Уровень логирования | `info` |

Полный список доступных параметров см. в `internal/config/config.go`.

## 🔬 Профилирование и оптимизация

Проект включает инструменты для профилирования CPU и памяти с использованием `net/http/pprof`.

### Профилирование CPU

```bash
# Запустить профиль на 30 секунд
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.pprof

# Анализ в терминале
go tool pprof cpu.pprof

# Или через веб-интерфейс
go tool pprof -http=:8081 cpu.pprof
```

### Профилирование памяти

```bash
curl http://localhost:8080/debug/pprof/heap > heap.pprof
go tool pprof -http=:8081 heap.pprof
```

### Доступные эндпоинты pprof

| Endpoint | Описание |
|----------|----------|
| `/debug/pprof/` | Индекс всех профилей |
| `/debug/pprof/heap` | Memory (Allocations) |
| `/debug/pprof/profile?seconds=30` | CPU (30 сек) |
| `/debug/pprof/goroutine` | Текущие горутины |
| `/debug/pprof/block` | Блокировки |

## 📈 История оптимизаций

### v1.0 - Исходная версия

**Проблемы:**
- Генератор ключей использовал `math/rand` с mutex - узкое место при параллельных запросах
- При создании URL без кастомного кода выполнялось до 5 SQL-запросов для проверки уникальности
- Validator создавался на каждый запрос (ненужные аллокации)
- Множественные вызовы `GetTraceID()` в каждом логе

**Результат профилирования:**
- `urlService.Create` занимал ~45% CPU времени
- До 5 SQL запросов на один Create-запрос

### v1.1 - Оптимизация генератора ключей

**Изменения:**
- Заменён `math/rand` на Snowflake-подобный алгоритм
- Гарантированная уникальность без проверки в БД
- Убран mutex в пользу атомарных операций

**Результат:**
- Генерация ключей: 119 ns/op → **31 ns/op** (3.8x быстрее)
- SQL запросов на Create: до 5 → **1**
- Время Create в профиле: 45% → ~29%

### v1.2 - Оптимизация middleware и handlers

**Изменения:**
- Переход с `gin.Default()` на `gin.New()` + ручное добавление middleware
- Validator перенесён в конструктор handler (один раз при старте)
- Убраны дублирующиеся вызовы trace_id в логах

**Результат:**
- Уменьшен overhead от framework
- Меньше аллокаций на каждый запрос

### Бенчмарки

```bash
go test -bench=. -benchmem ./internal/adapter/out/generator/
```

Результаты:
```
BenchmarkURLGenerator_Generate-12           37469672    31.92 ns/op    8 B/op    1 allocs/op
BenchmarkURLGenerator_Generate_Parallel-12   16590602    69.56 ns/op    8 B/op    1 allocs/op
```

## 🧪 Тестирование

### Unit-тесты

```bash
go test ./...
```

### Benchmark-тесты

```bash
go test -bench=. -benchmem ./internal/adapter/out/generator/
```

### Тест уникальности генератора

```bash
go test -v -run TestSnowflakeGenerator_Uniqueness ./internal/adapter/out/generator/
```
