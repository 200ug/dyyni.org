FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN apk add --no-cache gcc musl-dev

RUN CGO_ENABLED=1 go build -o blackbox_go .

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/blackbox_go /blackbox_go
EXPOSE 8080

CMD ["/blackbox_go"]

