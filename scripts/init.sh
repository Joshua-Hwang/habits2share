#!/bin/bash

dir=${1:-secrets}

mkdir -p $dir/

cat << EOF > .env
GOOGLE_WEB_CLIENT_ID=INSERT_YOUR_CLIENT_ID
GOOGLE_MOBILE_CLIENT_ID=INSERT_YOUR_CLIENT_ID
SESSIONS_FILE=$dir/sessions.csv
ACCOUNTS_FILE=$dir/accounts.json
HABITS_FILE=$dir/habits.json
TODO_FILE=$dir/todo.json
GOFLAGS=-tags=dev
EOF

cat << EOF > $dir/accounts.json
[
  {"Id": "testAccount1", "Email": "test1@mail.com"},
  {"Id": "testAccount2", "Email": "test2@mail.com"}
]
EOF
