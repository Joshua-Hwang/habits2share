#!/bin/bash

N=$1
if [[ -z "$N" ]]; then
  N=$(date +"%s")
fi

echo 'Using following value to shuffle test order'
echo "$N"
go test -cover -shuffle "$N" ./...
