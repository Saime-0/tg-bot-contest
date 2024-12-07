# Используем официальный образ Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY ../go.mod ../go.sum ./
RUN go mod download

# Копируем остальные файлы приложения
COPY --exclude=*.db .. .

# Собираем приложение
RUN go build -o ./bin/app ./cmd/app/main.go

# Используем более легкий образ для запуска
FROM alpine:latest

# Устанавливаем необходимые библиотеки
RUN apk --no-cache add sqlite sqlite-dev

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/contest-bot /app

# Указываем команду для запуска приложения
CMD ["./app"]
