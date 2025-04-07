#!/bin/bash

# Завершать скрипт при ошибках
set -e

DSN="${DSN:-}" # Строка соединения с БД
if [ -z "$DSN" ]; then
    echo "Ошибка: переменная окружения DSN не установлена."
    exit 1
fi

# Определить путь к директории с миграциями
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
MIGRATION_DIR=${SCRIPT_DIR}/migrations/main

# Получение названия последнего файла миграции
last_file=$(ls -1 $MIGRATION_DIR/*up.sql | tail -n 1)
if [ -z "$last_file" ]; then
    echo "Ошибка: файлы миграций не найдены в директории '$MIGRATION_DIR'."
    exit 1
fi

# Извлекаем имя файла без пути (например, "v1.0.0.up.sql")
last_version=$(basename "$last_file")

# Функция для получения текущей версии из БД
get_db_version() {
    sqlite3 "$DSN" "SELECT value FROM metadata WHERE key = 'db_version';" 2>/dev/null || echo ""
}

# Получаем текущую версию из БД
current_version=$(get_db_version)

# Проверяем, нужно ли выполнять инициализацию
if [ -z "$current_version" ] || [[ "$current_version" < "$last_version" ]]; then
    echo "Текущая версия ('$current_version') устарела или отсутствует. Выполняется инициализация до версии '$last_version'..."

    # Выполнение миграций
    {
        echo "BEGIN TRANSACTION;"
        cat "$MIGRATION_DIR"/*up.sql
        echo "INSERT OR REPLACE INTO metadata (key, value) VALUES ('db_version', '$last_version');"
        echo "COMMIT;"
    } | sqlite3 "$DSN"

    echo "Схема базы данных инициализирована до версии '$last_version'."
else
    echo "Текущая версия ('$current_version') актуальна. Инициализация не требуется."
fi