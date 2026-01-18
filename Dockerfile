# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot ./cmd/bot

# Runtime stage
FROM alpine:latest

# Метаданные образа
LABEL org.opencontainers.image.title="InPars Telegram Bot"
LABEL org.opencontainers.image.description="Telegram bot for monitoring real estate listings from InPars.ru API"
LABEL org.opencontainers.image.source="https://github.com/RedNessen/inpars_api_golang_bot_context_files"
LABEL org.opencontainers.image.licenses="MIT"

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 botuser && \
    adduser -D -u 1000 -G botuser botuser && \
    chown -R botuser:botuser /app

# Копируем собранное приложение
COPY --from=builder /app/bot .

# Переключаемся на непривилегированного пользователя
USER botuser

# Запускаем приложение
CMD ["./bot"]
