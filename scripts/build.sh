#!/usr/bin/env bash

[[ ! -f "package.json" ]] && echo "run from project root" && exit 1

mkdir -p go/bin
go build -C go -o bin/blackbox_go .
bun run build

