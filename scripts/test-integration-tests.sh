#!/bin/bash

export SESSIONS_FILE=secrets_integration/sessions.csv
export ACCOUNTS_FILE=secrets_integration/accounts.json
export HABITS_FILE=secrets_integration/habits.json
export TODO_FILE=secrets_integration/todo.json
export GOFLAGS=-tags=dev

#./scripts/build-frontend.sh || exit $?
./scripts/build-backend.sh || exit $?

trap 'pids=( $(jobs -p) ); [ -n "$pids" ] && kill -- "${pids[@]/#/-}"' EXIT

restart() {
  pids=( $(jobs -p) )
  [ -n "$pids" ] && kill -- "${pids[@]/#/-}"
  rm -rf secrets_integration/
  ./scripts/init.sh secrets_integration
  setsid ./build/server &
  until curl -s --head localhost:8080/healthcheck > /dev/null; do sleep 1; done
}

restart

tmp_file1=$(mktemp)
curl -s -F email=test1@mail.com -c $tmp_file1 localhost:8080/login

res=$(curl -s -b $tmp_file1 localhost:8080/my/habits)
if [[ "$res" != '[]' ]]; then
  echo "Failed test at line $LINENO"
  echo "Received"
  echo "$res"
  exit 1
fi

curl -s -X POST -c $tmp_file1 -b $tmp_file1 localhost:8080/logout

res=$(curl -s -b $tmp_file1 localhost:8080/my/habits)
if [[ "$res" != 'Anonymous access forbidden' ]]; then
  echo "Failed test at line $LINENO"
  echo "Received"
  echo "$res"
  exit 1
fi

restart

tmp_file1=$(mktemp)
curl -s -F email=test1@mail.com -c $tmp_file1 localhost:8080/login
tmp_file2=$(mktemp)
curl -s -F email=test2@mail.com -c $tmp_file2 localhost:8080/login

./scripts/generate-data/share_habit.sh 2> /dev/null

res=$(curl -s -b $tmp_file1 localhost:8080/my/habits)
if (( $(echo "$res" | jq 'length') != 1 )); then
  echo "Failed test at line $LINENO"
  echo "Received"
  # echo "$res" | jq
  exit 1
fi
res=$(curl -s -b $tmp_file2 localhost:8080/shared/habits)
if (( $(echo "$res" | jq 'length') != 1 )); then
  echo "Failed test at line $LINENO"
  echo "Received"
  echo "$res" | jq
  exit 1
fi

restart

tmp_file1=$(mktemp)
curl -s -F email=test1@mail.com -c $tmp_file1 localhost:8080/login
tmp_file2=$(mktemp)
curl -s -F email=test2@mail.com -c $tmp_file2 localhost:8080/login

./scripts/generate-data/100_habits_each.sh 2> /dev/null

res=$(curl -s -b $tmp_file1 localhost:8080/my/habits)
if (( $(echo "$res" | jq 'length') != 100 )); then
  echo "Failed test at line $LINENO"
  echo "Received"
  echo "$res" | jq
  exit 1
fi
res=$(curl -s -b $tmp_file2 localhost:8080/my/habits)
if (( $(echo "$res" | jq 'length') != 100 )); then
  echo "Failed test at line $LINENO"
  echo "Received"
  echo "$res" | jq
  exit 1
fi
