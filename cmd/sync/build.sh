#!/bin/sh

source .env
(cd $GOPATH/src/google.golang.org/api/photoslibrary/v1 && git apply $CI_PROJECT_DIR/patches/media-item-filename.patch)

set -e 
FLAGS="-ldflags \"-X main.GoogleClientId=$GOOGLE_CLIENT_ID -X main.GoogleClientSecret=$GOOGLE_CLIENT_SECRET\""

mkdir -p builds
go get
GOOS=linux   GOARCH=386   go build $FLAGS -o builds/gphotosync_linux_386
GOOS=linux   GOARCH=amd64 go build $FLAGS -o builds/gphotosync_linux_amd64
GOOS=linux   GOARCH=arm   go build $FLAGS -o builds/gphotosync_linux_arm7
GOOS=darwin  GOARCH=amd64 go build $FLAGS -o builds/gphotosync_mac_amd64
GOOS=windows GOARCH=386   go build $FLAGS -o builds/gphotosync_windows_386.exe
GOOS=windows GOARCH=amd64 go build $FLAGS -o builds/gphotosync_windows_amd64.exe