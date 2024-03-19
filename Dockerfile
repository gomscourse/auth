FROM golang:1.22.1-alpine3.19 AS builder

COPY . /github.com/gomscourse/auth/source/
WORKDIR /github.com/gomscourse/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/main.go

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

WORKDIR /root/
COPY --from=builder /github.com/gomscourse/auth/source/bin/auth_server .
COPY --from=builder /github.com/gomscourse/auth/source/entrypoint.sh .
COPY --from=builder /github.com/gomscourse/auth/source/migrations ./migrations

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose