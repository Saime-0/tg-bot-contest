#!/bin/bash

# Завершать скрипт при ошибках
set -e

# Получение пути до директории с миграциями и названия файла БД из переменных окружения
MIGRATION_DIR="${MIGRATION_DIR:-}" # Путь до директории с миграциями
DB_FILE="${DB_FILE:-}" # Название файла базы данных

# Проверка на наличие обязательных переменных окружения
if [ -z "$MIGRATION_DIR" ]; then
    echo "Ошибка: переменная окружения MIGRATION_DIR не установлена."
    exit 1
fi

if [ -z "$DB_FILE" ]; then
    echo "Ошибка: переменная окружения DB_FILE не установлена."
    exit 1
fi

# Получение названия последнего файла миграции
last_file=$(ls -1 "$MIGRATION_DIR"/*up.sql | tail -n 1)
db_version=$(basename "$last_file")  # Получаем имя файла без пути

# Выполнение миграций
echo "Выполнение миграций..."
{
    echo "BEGIN TRANSACTION;"
    cat "$MIGRATION_DIR"/*up.sql
    echo "INSERT OR REPLACE INTO metadata (key, value) VALUES ('db_version', '$db_version');"
    echo "COMMIT;"
} | sqlite3 "$DB_FILE"

echo "База данных инициализирована"
