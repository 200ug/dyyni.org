#!/usr/bin/env bash

set -euo pipefail

if [ $# -ne 1 ]; then
    echo "usage: $0 username@host"
    exit 1
fi

[[ ! -f ".env" ]] && echo "[!] .env missing, remember to add it to server before deploying"

DEST="$1"
REMOTE_DIR="~/services/dyyni_blackbox"

rsync -avz --relative --ignore-missing-args --exclude='bin/' \
    server/ \
    .env \
    scripts/deploy_backend.sh \
    "$DEST:$REMOTE_DIR/"
