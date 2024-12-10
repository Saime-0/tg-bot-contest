#!/bin/bash

# Завершать скрипт при ошибках
set -e

# Получение пути до директории с миграциями и названия файла БД из переменных окружения
MIGRATION_DIR="${MIGRATION_DIR:-}" # Путь до директории с миграциями
if [ -z "$MIGRATION_DIR" ]; then
    echo "Ошибка: переменная окружения MIGRATION_DIR не установлена."
    exit 1
fi

DSN="${DSN:-}" # Строка соединения с БД
if [ -z "$DSN" ]; then
    echo "Ошибка: переменная окружения DSN не установлена."
    exit 1
fi


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
