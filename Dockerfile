FROM golang:1.24.1-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем go.mod и go.sum отдельно для кэша
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app main.go


# Runtime stage

FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/app /app/app

# Открываем порт
EXPOSE 8080

# Запуск
ENTRYPOINT ["/app/app"]