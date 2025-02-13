FROM golang:1.22 as builder

WORKDIR /app

# Копируем файлы в контейнер
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN go build -o wallet_go ./cmd/main.go

# Финальный образ
FROM debian:bookworm

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/wallet_go .

# Копируем конфиг окружения
COPY config.env .

# Запуск приложения
CMD ["./wallet_go"]
