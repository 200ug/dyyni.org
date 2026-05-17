#!/usr/bin/env bash

[[ ! -f "package.json" ]] && echo "run from project root" && exit 1

mkdir -p database
podman compose up -d --build

