#!/bin/sh

# Load oauthClientIds in .env
export $(cat .env | xargs)

go run ./cmd/http
