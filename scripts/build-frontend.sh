#!/bin/bash

set -euo pipefail

cd frontend
yarn install
yarn build
cd ..
