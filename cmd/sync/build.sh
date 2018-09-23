#!/bin/bash

source ./.env

set -e
build() {
    go build -ldflags "-X main.GoogleClientId=$GOOGLE_CLIENT_ID -X main.GoogleClientSecret=$GOOGLE_CLIENT_SECRET" -o builds/gphotosync-$GOOS-$GOARCH$BINEXT
}

mkdir -p builds
go get
GOOS=linux   GOARCH=386 build
GOOS=linux   GOARCH=amd64 build
GOOS=linux   GOARCH=arm   build
GOOS=darwin  GOARCH=amd64 build
BINEXT=".exe" GOOS=windows GOARCH=386   build
BINEXT=".exe" GOOS=windows GOARCH=amd64 build