#!/bin/bash

tmp_file1=$(mktemp)

#echo $tmp_file1

curl -F email=test1@mail.com -c $tmp_file1 localhost:8080/login
curl -b $tmp_file1 --json "{\"Name\":\"simple\", \"Description\": \"simple description\", \"Frequency\": 5}" localhost:8080/my/habits
