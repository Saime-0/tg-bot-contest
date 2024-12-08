# Используем официальный образ Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /usr/src/app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.* .
RUN go mod download && go mod verify

# Копируем остальные файлы приложения
COPY . .

# Собираем приложение
RUN go build -v -o /usr/local/bin/app ./cmd/app/main.go

# Используем более легкий образ для запуска
FROM alpine:latest

# Устанавливаем необходимые библиотеки
RUN apk --no-cache add sqlite sqlite-dev # libc6-compat

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

# Указываем команду для запуска приложения
CMD ["app"]
#CMD ["tail", "-f", "/dev/null"]
