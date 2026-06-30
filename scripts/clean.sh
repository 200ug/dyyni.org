#!/usr/bin/env bash

[[ ! -d ".git" ]] && echo "run from project root" && exit 1

# build artifacts
rm -rf server/bin
rm -rf dist

