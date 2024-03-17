FROM golang:1.22.1-alpine3.19 AS builder

COPY . /github.com/gomscourse/auth/source/
WORKDIR /github.com/gomscourse/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/gomscourse/auth/source/bin/auth_server .
COPY --from=builder /github.com/gomscourse/auth/source/entrypoint.sh .