#!/bin/bash

# Завершать скрипт при ошибках
set -e

DSN="${DSN:-}" # Строка соединения с БД
if [ -z "$DSN" ]; then
    echo "Ошибка: переменная окружения DSN не установлена."
    exit 1
fi

MIGRATION_DIR=migrations/main # Путь до директории с миграциями

last_file=$(ls -1 "$MIGRATION_DIR"/*up.sql | tail -n 1) # Получение названия последнего файла миграции
db_version=$(basename "$last_file")  # Получаем имя файла без пути

# Выполнение миграций
echo "Выполнение инициализации..."
{
    echo "BEGIN TRANSACTION;"
    cat "$MIGRATION_DIR"/*up.sql
    echo "INSERT OR REPLACE INTO metadata (key, value) VALUES ('db_version', '$db_version');"
    echo "COMMIT;"
} | sqlite3 "$DSN"

echo "Схема база данных инициализирована"
