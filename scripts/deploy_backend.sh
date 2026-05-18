#!/usr/bin/env bash

[[ ! -f "go/Dockerfile" ]] && echo "run from project root" && exit 1

mkdir -p database
podman rm -f bb-backend 2>/dev/null
podman build -t blackbox_go go/
podman run -d \
    --name bb-backend \
    --network host \
    --restart unless-stopped \
    -e ENV=production \
    -v ./database:/database \
    --health-cmd "wget -q -O /dev/null http://127.0.0.1:8081/health" \
    --health-interval 10s \
    --health-timeout 5s \
    --health-retries 3 \
    --health-start-period 5s \
    blackbox_go