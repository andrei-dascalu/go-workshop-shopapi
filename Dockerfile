FROM golang:1.15-alpine AS base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED=0

# System dependencies
RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  git \
  && update-ca-certificates

# Application dependencies
COPY . .
RUN go mod download \
  && go mod verify

### Development
FROM base AS dev
WORKDIR /app

# Hot reloading mod
RUN go get github.com/githubnemo/CompileDaemon
EXPOSE 8080

ENTRYPOINT ["CompileDaemon","-directory=/app", "-command", "/app/go-workshop-shopapi"]

### Executable builder
FROM base AS builder
WORKDIR /app

RUN go build -o ./go-workshop-shopapi -a .

### Production
FROM alpine:latest

ARG APP_VERSION="v0.0.1"
ENV APP_VERSION "${APP_VERSION}"

RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  git \
  tzdata \
  && update-ca-certificates

# Copy executable and use an unprivileged user
COPY --from=builder /app/go-workshop-shopapi /go/bin/go-workshop-shopapi
EXPOSE 8080

ENTRYPOINT ["/go/bin/go-workshop-shopapi"]
