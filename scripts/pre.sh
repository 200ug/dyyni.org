#!/usr/bin/env bash

bun run format
bun run lint

go mod tidy -C go
gofmt -w go

