#!/usr/bin/env bash

[[ ! -f "package.json" ]] && echo "run from project root" && exit 1

# build artifacts
rm -rf go/bin
rm -rf dist

