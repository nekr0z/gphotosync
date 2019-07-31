#!/bin/sh
source ./.env
go build -ldflags "-X main.GoogleClientId=$GOOGLE_CLIENT_ID -X main.GoogleClientSecret=$GOOGLE_CLIENT_SECRET" -v
