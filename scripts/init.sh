#!/bin/bash

mkdir -p secrets/

cat << EOF > .env
GOOGLE_WEB_CLIENT_ID=INSERT_YOUR_CLIENT_ID
GOOGLE_MOBILE_CLIENT_ID=INSERT_YOUR_CLIENT_ID
SESSIONS_FILE=secrets/sessions.csv
ACCOUNTS_FILE=secrets/accounts.json
HABITS_FILE=secrets/habits.json
TODO_FILE=secrets/todo.json
GOFLAGS=-tags=dev
EOF

cat << EOF > secrets/accounts.json
[
  {"Id": "testAccount1", "Email": "test1@mail.com"},
  {"Id": "testAccount2", "Email": "test2@mail.com"}
]
EOF
