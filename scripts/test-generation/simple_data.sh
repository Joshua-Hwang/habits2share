#!/bin/bash

tmp_file1=$(mktemp)

#echo $tmp_file1

curl -F email=test1@mail.com -c $tmp_file1 localhost:8080/login

for ((i=0;i<100;i++)); do
  name=$(uuidgen)
  description="$name description has $(uuidgen)"
  frequency=$(( RANDOM % 7 + 1 ))
  curl -b $tmp_file1 --json "{\"Name\":\"$name\", \"Description\": \"$description\", \"Frequency\": 3}" localhost:8080/my/habits
  echo
done

tmp_file2=$(mktemp)

#echo $tmp_file2

curl -F email=test2@mail.com -c $tmp_file2 localhost:8080/login

for ((i=0;i<100;i++)); do
  name=$(uuidgen)
  description="$name description has $(uuidgen)"
  frequency=$(( RANDOM % 7 + 1 ))
  curl -b $tmp_file2 --json "{\"Name\":\"$name\", \"Description\": \"$description\", \"Frequency\": 3}" localhost:8080/my/habits
  echo
done
