FROM golang:1.14-alpine AS builder

ENV \
  GOOS=linux \
  GOARCH=amd64

RUN apk add --no-cache git
RUN go get \
  golang.org/x/crypto/nacl/box \
  github.com/google/subcommands

COPY /src /go/src/nacl-cli
WORKDIR /go/src/nacl-cli

RUN go build -o /go/bin/nacl .



FROM scratch

COPY --from=builder /go/bin/nacl .
