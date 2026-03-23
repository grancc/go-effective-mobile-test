# Сборка приложения
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /gotodo ./cmd/main.go

# Финальный образ
FROM alpine:3.19

WORKDIR /app

# Копируем бинарник и конфиги
COPY --from=builder /gotodo .
COPY configs ./configs

EXPOSE 8080

CMD ["./gotodo"]
