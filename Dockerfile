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

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/bot .

# Запускаем приложение
CMD ["./bot"]
