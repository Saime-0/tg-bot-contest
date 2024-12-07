# tg-bot-contest

Перед первым деплоем необходимо:

- Скопировать проект на удаленный сервер и подключиться по ssh
- Установить `docker`, `sqlite`
- Перейти в директорию проекта
- Выполнить скрипт для инициализации БД

```sh
mkdir -p ~/.local/share/contest-bot
MIGRATION_DIR="./migrations/main" DB_FILE="~/.local/share/contest-bot/main.db" ./scripts/init-scheme.sh
```
