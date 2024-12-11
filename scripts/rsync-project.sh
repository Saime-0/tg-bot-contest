#!/bin/sh

rsync -avzrP --include='cmd/***' \
          --include='deploy/***' \
          --include='internal/***' \
          --include='migrations/***' \
          --include='scripts/***' \
          --include='go.mod' \
          --include='go.sum' \
          --exclude='*' \
          . ${SSH_HOST}:${DIR}