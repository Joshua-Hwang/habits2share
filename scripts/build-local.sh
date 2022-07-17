#!/bin/sh

cd frontend
yarn install
yarn build
cd ..

./scripts/build.sh
