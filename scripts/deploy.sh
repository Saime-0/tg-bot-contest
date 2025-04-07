#!/bin/sh

set -e

# Имя хоста из .ssh/config
ssh_host=$SSH_HOST
if [ -z "$SSH_HOST" ]; then
    echo "Ошибка: переменная окружения SSH_HOST не установлена. SSH_HOST=vds-server-a"
    exit 1
fi

# Версия приложения, понадобится для именования контейнеров
version=$VERSION
if [ -z "$VERSION" ]; then
    echo "Ошибка: переменная окружения VERSION не установлена. Например VERSION=v1.0.1"
    exit 1
fi

# Префикс контейнеров, нужен для возможности одновременно запускать несколько ботов, но для разных нужд
app_prefix=$APP_PREFIX
if [ -z "$APP_PREFIX" ]; then
    echo "Ошибка: переменная окружения APP_PREFIX не установлена. Например APP_PREFIX=dev"
    exit 1
fi

# Токен бота
if [ -z "$TOKEN" ]; then
    echo "Ошибка: переменная окружения TOKEN не установлена. Можно взять у t.me/BotFather"
    exit 1
fi

# Директория для загрузки исходного кода приложения
src_dir="~/src/${app_prefix}-tg-contest-bot"
# Рабочая директория приложения
release_dir="~/release/${app_prefix}-tg-contest-bot"

# Проверить наличие требуемых утилит на сервере
#ssh $ssh_host "bash -s" < "$(cd "$(dirname "$0")" && pwd)"/check-utils.sh
ssh $ssh_host "bash -s" < ./scripts/check-utils.sh

# Создать директорию для исходного кода
ssh $ssh_host "mkdir -p $src_dir"
# Скопировать исходный код на хост
SSH_HOST=$ssh_host DIR="$src_dir" ./scripts/rsync-project.sh


# Создать директорию с релизом
ssh $ssh_host "mkdir -p $release_dir"
# Выполнить скрипт для инициализации БД
database_dir=$release_dir
database_filename="main.db"
ssh $ssh_host "DSN=file:$database_dir/$database_filename $src_dir/scripts/init-scheme.sh"

# Создать образ
image_base_name="${app_prefix}-tg-contest-bot"
new_image_name="${image_base_name}"
image_script="$src_dir/deploy/Dockerfile"
ssh $ssh_host "podman build --tag $new_image_name:$version -f $image_script $src_dir"

echo "Поиск прошлого контейнера ..."
previous_container_id=$(ssh $ssh_host "podman ps -q --no-trunc --format \{\{.ID\}\} --filter \"name=${image_base_name}\"")
if [ -n "$previous_container_id" ]; then
  previous_container_name=$(ssh $ssh_host "podman ps -q --format \{\{.Names\}\} --filter \"id=${previous_container_id}\"")
  echo "Найден прошлый контейнер с ID $previous_container_id и именем $previous_container_name"
  ssh $ssh_host "podman stop $previous_container_id" # Остановить прошлый контейнер
  echo "Прошлый контейнер остановлен"
else
  echo "Прошлых контейнеров нет"
fi

# Если произойдет ошибка, работа скрипта будет продолжена
set +e

# Запустить новый контейнер
echo "Запуск нового контейнера ..."
new_container_name=${new_image_name}_$(date +"%Y%m%d_%H%M%S")
db_dsn="file:/var/lib/tg-contest-bot/$database_filename"
ssh $ssh_host "podman run --name "$new_container_name" \
  --cpus=0.4 \
  --memory=256m \
  --replace \
  --detach \
  -e TOKEN=$TOKEN \
  -e MAIN_DATABASE_DSN=$db_dsn \
  --restart unless-stopped \
  --log-driver=journald \
  -v $release_dir:/var/lib/tg-contest-bot \
  $new_image_name:$version"

if [ $? -eq 0 ]; then
  if [ -n "$previous_container_id" ]; then
    echo "Удаление контейнера ID $previous_container_id"
    ssh $ssh_host "podman rm $previous_container_id"
  fi
  echo ""
  echo "====================================="
  echo "Запущен контейнер новый контейнер с приложением."
  echo "Новый контейнер: $new_container_name"
  echo "На основе образа: $new_image_name"
  echo "Версия приложения: $version"
  echo "Префикс приложения: $app_prefix"
  echo "Прошлый остановленный и удаленный контейнер: ${previous_container_name:-Нет}"
  echo "Прошлый остановленный и удаленный контейнер: ${previous_container_id:-Нет}"
  echo "====================================="
  exit 0;
else
  echo ""
  echo "====================================="
  echo "Не удалось запустить контейнер с новой версией приложения"
fi


# Если произойдет ошибка, скрипт остановится
set -e

# Запустить прошлый контейнер, если запуск нового завершился с ошибкой
if [ -n "$previous_container_id" ]; then
  echo "Начат запуск прошлого контейнера"
  ssh "$ssh_host" "podman start $previous_container_id"
  echo "Прошлый контейнер успешно запущен"
fi

