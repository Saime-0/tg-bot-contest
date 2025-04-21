#!/bin/bash

set -e

for util in podman sqlite3 rsync cpulimit jq goose; do
  command -v "$util" &>/dev/null || { echo "Ошибка: Утилита '$util' не найдена."; exit 1; }
done

echo "Все требуемые утилиты установлены на хосте $HOST"
exit 0