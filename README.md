# tg-bot-contest

Перед первым деплоем необходимо:

- Установить `docker`, `sqlite`, `rsync`
- Скопировать проект на удаленный сервер

```shell
SSH_HOST=server-ssh-host DIR="~/src/tgbotcontest/" ./scripts/rsync-project.sh
```

- Подключиться по ssh и перейти в директорию проекта

```shell
ssh server-ssh-host
cd ~/src/tgbotcontest
```

- Выполнить скрипт для инициализации БД

```shell
mkdir -p ~/.local/share/contest-bot
DSN=file:~/.local/share/contest-bot/main.db ./scripts/init-scheme.sh
```