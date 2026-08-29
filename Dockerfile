# Сборка программы
FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY /backend/go.mod /backend/go.sum ./
RUN go mod download
COPY backend/ .

# Для SQLite важна переменная CGO_ENABLED=1, так как драйвер SQLite использует C-библиотеку
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Запуск программы
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
