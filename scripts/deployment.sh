#!/bin/bash

print_help() {
  # Описание переменных
  echo "Используемые переменные:"
  echo "- NAMESPACE: пространство имен (напр. dev, prod, creeper007)"
  echo "- APP_SRC_DIR: путь к директории с приложения"
  echo "- APP_FILENAME: имя исполняемого файла приложения"
  echo "- VERSION: версия приложения, например, ветка или тег Git"
  echo "- TOKEN: токен для доступа к TelegramBotAPI"
}

# Проверка наличия обязательных переменных окружения
check_variable() {
    if [ -z "${!1}" ]; then
        echo "Ошибка: Переменная '$1' не задана."
        print_help
        exit 1
    fi
}

# Обязательные переменные
check_variable NAMESPACE
check_variable APP_SRC_DIR
check_variable APP_FILENAME
check_variable VERSION
check_variable TOKEN

#set -o xtrace

set -e

QUADLETS_DIR=$HOME/.config/containers/systemd
mkdir -p "$QUADLETS_DIR"

# Создать quadlet volume из шаблона, если его нет, и запустить службу
VOLUME_NAME="${NAMESPACE}-tg-contest-bot"
if ! podman volume exists "$VOLUME_NAME"; then
  QUADLET_TPL_PATH="$APP_SRC_DIR/scripts/quadlet/volume.template"
  QUADLET_PATH="$QUADLETS_DIR/$VOLUME_NAME.volume"
  cp "$QUADLET_TPL_PATH" "$QUADLET_PATH" #  envsubst < $QUADLET_TPL_PATH > "$QUADLET_PATH"
  systemctl --user daemon-reload
  systemctl --user start "$VOLUME_NAME-volume"
fi


# Создать БД если не создана
DB_DIR=$(podman volume inspect "systemd-$VOLUME_NAME" | jq -r '.[0].Mountpoint')
DB_FILENAME=main.db
DB_PATH="$DB_DIR/$DB_FILENAME"
if [ ! -f "$DB_PATH" ]; then
  sqlite3 "$DB_PATH" ''
  goose --dir="$APP_SRC_DIR/migrations/main" sqlite3 "$DB_PATH" up
fi


# Создает quadlet container из шаблона и запустить службу
CONTAINER_NAME="${NAMESPACE}-tg-contest-bot"
APP_PATH="$APP_SRC_DIR/bin/$APP_FILENAME"
QUADLET_TPL_PATH="$APP_SRC_DIR/scripts/quadlet/container.template"
QUADLET_PATH="$QUADLETS_DIR/$CONTAINER_NAME.container"
export CONTAINER_NAME TOKEN DB_FILENAME VOLUME_NAME APP_PATH APP_FILENAME VERSION
envsubst < "$QUADLET_TPL_PATH" > "$QUADLET_PATH"
systemctl --user daemon-reload
systemctl --user restart "$CONTAINER_NAME"



