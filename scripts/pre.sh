#!/usr/bin/env bash

if [[ "$1" != "go" ]]; then
    bun run format
    bun run lint
fi

if [[ "$1" != "bun" ]]; then
    go mod tidy -C go
    gofmt -w go
fi

