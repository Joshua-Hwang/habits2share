#!/bin/sh

# Load GOFLAGS
export $(cat .env | xargs)

# go run is out of the question because debug info
mkdir -p build
go build -o ./build/server ./cmd/http && ./build/server
