#!/bin/bash
# Code generated by DeepSeek

# Завершать скрипт при ошибках и включить дополнительные проверки
set -euo pipefail

# Проверка и настройка переменных окружения
DSN="${DSN:-}"
if [[ -z "$DSN" ]]; then
    echo "Ошибка: переменная окружения DSN не установлена." >&2
    exit 1
fi

# Определение путей
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
MIGRATION_DIR="${SCRIPT_DIR}/migrations/main"

# Проверка существования директории с миграциями
if [[ ! -d "$MIGRATION_DIR" ]]; then
    echo "Ошибка: директория миграций '$MIGRATION_DIR' не существует." >&2
    exit 1
fi

# Получение и сортировка файлов миграций
migrations=()
while IFS= read -r -d $'\0' file; do
    migrations+=("$file")
done < <(find "$MIGRATION_DIR" -maxdepth 1 -name '*up.sql' -print0 | sort -z)

if [[ ${#migrations[@]} -eq 0 ]]; then
    echo "Ошибка: файлы миграций не найдены в директории '$MIGRATION_DIR'." >&2
    exit 1
fi

# Получение последней версии
last_version=$(basename "${migrations[-1]}")

# Функция для получения текущей версии из БД
get_db_version() {
    local version
    version=$(sqlite3 "$DSN" "SELECT value FROM metadata WHERE key = 'db_version';" 2>/dev/null || echo "")
    echo "$version"
}

# Получаем текущую версию из БД
current_version=$(get_db_version)

# Если БД не инициализирована, создаем таблицу metadata
if [[ -z "$current_version" ]]; then
    echo "Инициализация новой базы данных..."
    sqlite3 "$DSN" "CREATE TABLE IF NOT EXISTS metadata (key TEXT PRIMARY KEY, value TEXT);"
    current_version=""
fi

# Проверка необходимости выполнения миграций
if [[ -z "$current_version" ]]; then
    echo "База данных пустая. Применяем все миграции..."
    pending_migrations=("${migrations[@]}")
elif [[ "$current_version" == "$last_version" ]]; then
    echo "Текущая версия ('$current_version') актуальна. Миграции не требуются."
    exit 0
else
    echo "Текущая версия ('$current_version') устарела. Определяем необходимые миграции..."

    # Находим индекс текущей версии
    current_index=-1
    for i in "${!migrations[@]}"; do
        if [[ $(basename "${migrations[i]}") == "$current_version" ]]; then
            current_index=$i
            break
        fi
    done

    if [[ $current_index -eq -1 ]]; then
        echo "Ошибка: текущая версия '$current_version' не найдена в списке миграций." >&2
        exit 1
    fi

    # Получаем миграции для применения
    pending_migrations=("${migrations[@]:$((current_index + 1))}")
fi

# Применяем миграции, если они есть
if [[ ${#pending_migrations[@]} -gt 0 ]]; then
    echo "Применяем следующие миграции:"
    printf ' - %s\n' "${pending_migrations[@]##*/}"

    # Собираем SQL для выполнения
    sql_script=$(mktemp)
    trap 'rm -f "$sql_script"' EXIT

    {
        echo "BEGIN TRANSACTION;"
        for file in "${pending_migrations[@]}"; do
            echo "-- Applying: $(basename "$file")"
            cat "$file"
            echo ""
        done
        echo "INSERT OR REPLACE INTO metadata (key, value) VALUES ('db_version', '$last_version');"
        echo "COMMIT;"
    } > "$sql_script"

    # Выполняем миграции
    if ! sqlite3 "$DSN" < "$sql_script"; then
        echo "Ошибка при выполнении миграций!" >&2
        exit 1
    fi

    echo "Миграции успешно применены. Текущая версия: '$last_version'"
else
    echo "Нет новых миграций для применения."
fi