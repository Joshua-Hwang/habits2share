#!/bin/sh

./scripts/build-frontend.sh

docker build -t habits2share .
