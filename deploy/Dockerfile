# Используем официальный образ Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /src

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.* .
RUN go mod download && go mod verify

# Копируем остальные файлы приложения
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 go build -v -o /bin/app ./cmd/app/main.go

# Используем более легкий образ для запуска
FROM alpine:latest

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /bin/app /usr/local/bin/app

# Указываем команду для запуска приложения
CMD ["app"]
