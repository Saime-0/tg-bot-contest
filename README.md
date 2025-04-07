# tg-bot-contest

### ci

[![prod](https://github.com/Saime-0/tg-bot-contest/actions/workflows/deploy-prod.yml/badge.svg)](https://github.com/Saime-0/tg-bot-contest/actions/workflows/deploy-prod.yml)

[![staging](https://github.com/Saime-0/tg-bot-contest/actions/workflows/deploy-staging.yml/badge.svg)](https://github.com/Saime-0/tg-bot-contest/actions/workflows/deploy-staging.yml)

### about

Перед первым деплоем необходимо:

- Установить `docker`, `sqlite`, `rsync` на сервере
- Скопировать проект на удаленный сервер

```shell
ssh server-ssh-host "mkdir -p ~/src/tgbotcontest/"
SSH_HOST=server-ssh-host DIR="~/src/tgbotcontest/" ./scripts/rsync-project.sh
```

- Подключиться по ssh и перейти в директорию проекта

```shell
ssh server-ssh-host
cd ~/src/tgbotcontest
```

- Выполнить скрипт для инициализации БД

```shell
mkdir -p ~/opt/contest-bot
DSN=file:~/opt/contest-bot/main.db ./scripts/init-scheme.sh
```