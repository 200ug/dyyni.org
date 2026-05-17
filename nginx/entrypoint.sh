#!/bin/sh

CERT_DIR=/etc/nginx/certs

mkdir -p "$CERT_DIR"

if [ ! -f "$CERT_DIR/cert.pem" ]; then
    openssl req -x509 -nodes -days 3650 \
        -newkey rsa:2048 \
        -keyout "$CERT_DIR/key.pem" \
        -out "$CERT_DIR/cert.pem" \
        -subj "/CN=organic.dyyni.org"
fi

exec nginx -g 'daemon off;'

