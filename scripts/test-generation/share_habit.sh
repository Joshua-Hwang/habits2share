#!/bin/bash

tmp_file1=$(mktemp)

#echo $tmp_file1

curl -F email=test1@mail.com -c $tmp_file1 localhost:8080/login
habit_id=$(curl -b $tmp_file1 --json "{\"Name\":\"shared\", \"Description\": \"sharing with test2\", \"Frequency\": 3}" localhost:8080/my/habits)

curl -b $tmp_file1 -X POST localhost:8080/user/testAccount2/habit/$habit_id
