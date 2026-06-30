#!/usr/bin/env bash

[[ ! -d ".git" ]] && echo "run from project root" && exit 1

mkdir -p server/bin
go build -C server -o bin/blackbox .
bun run build

