# --- Базовая конфигурация ---
ARG GO_VERSION=1.25

# --- Builder stage (prod) ---
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем статические конфиги
COPY config.yaml .env ./

# Копируем всё остальное (только для сборки)
COPY . .

# Сборка (статическая, без cgo)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main ./cmd/app

# --- Минимальный prod-образ ---
FROM alpine AS prod

# Установка сертификатов и создание non-root пользователя
RUN apk --no-cache add ca-certificates \
    && addgroup -g 1001 -S appgroup \
    && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Копируем только нужные файлы
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /app/config.yaml .

# Применяем non-root пользователя
USER appuser

EXPOSE 8080
CMD ["./main"]

# --- Debug-образ (отдельный слой, без дублирования исходников) ---
# Можно использовать как `--target=debug`
FROM builder AS debug

RUN apk --no-cache add delve

# Пересобираем с отладочной информацией
RUN go build \
    -gcflags="all=-N -l" \
    -o /app/main_debug ./cmd/app

EXPOSE 8080 40000
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/main_debug"]
