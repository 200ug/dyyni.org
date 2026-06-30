#!/usr/bin/env bash

if [[ "$1" != "go" ]]; then
    bun run --cwd web format
    bun run --cwd web lint
fi

if [[ "$1" != "bun" ]]; then
    go mod tidy -C server
    gofmt -w server
fi

