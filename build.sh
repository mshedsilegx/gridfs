#!/bin/bash

# Syntax: build.sh 0.5.0xg-20241024

APP_VERSION="$1"

mkdir -p bin
[ ! -s go.mod ] && go mod init mckesson/gridfs
go mod tidy
go build -v -ldflags "-s -w -X main.Version=${APP_VERSION}" -o ./bin/

tar Jcvf ../gridfs-${APP_VERSION%xg*}-amd64.tar.xz *
