#!/usr/bin/env bash

[[ ! -f "server/Dockerfile" ]] && echo "run from project root" && exit 1

podman rm -f blackbox 2>/dev/null
podman build -t blackbox server/
podman run -d \
    --name blackbox \
    --network host \
    --restart unless-stopped \
    -e ENV=production \
    -v ./.env:/.env \
    --health-cmd "wget -q -O /dev/null http://127.0.0.1:8081/health" \
    --health-interval 10s \
    --health-timeout 5s \
    --health-retries 3 \
    --health-start-period 5s \
    blackbox
