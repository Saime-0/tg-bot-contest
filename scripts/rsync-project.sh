#!/bin/sh

rsync --delete -avzrP \
    --include='bin/***' \
    --include='migrations/***' \
    --include='scripts/***' \
    --exclude='*' \
    . "$SSH_HOST:$DIR"