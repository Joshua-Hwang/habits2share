#!/bin/sh

# Load GOFLAGS
export $(cat .env | xargs)

# go run is out of the question because debug info
./scripts/build-backend.sh && ./build/server
