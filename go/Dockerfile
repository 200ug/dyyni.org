FROM docker.io/library/golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o blackbox_go .

FROM docker.io/library/alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/blackbox_go /blackbox_go
CMD ["/blackbox_go"]

