#!/bin/sh

mkdir -p build
go build -o ./build/server ./cmd/http

cd frontend
yarn install
yarn build
cd ..
