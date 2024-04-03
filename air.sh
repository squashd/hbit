#!/bin/zsh

BASE_PATH=$(pwd)

start_air() {
  cd "$BASE_PATH/$1" || exit
  air &
}

start_air "./cmd/api-gateway"
start_air "./cmd/task-service"
start_air "./cmd/user-service"
start_air "./cmd/rpg-service"
start_air "./cmd/feat-service"
start_air "./cmd/update-service"

wait
