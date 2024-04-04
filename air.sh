#!/bin/zsh

BASE_PATH=$(pwd)

free_port() {
  local port=$1
  echo "Attempting to free up port $port..."
  local pid=$(lsof -ti tcp:$port)
  if [[ -n $pid ]]; then
    echo "Killing process $pid on port $port"
    kill -9 $pid
  fi
}

start_air() {
  service_path=$1
  port=$2
  free_port $port
  cd "$BASE_PATH/$service_path" || exit
  echo "Starting service in $service_path on port $port"
  air &
}

start_air "./cmd/api-gateway" 9080
start_air "./cmd/task-service" 9081
start_air "./cmd/user-service" 9082
start_air "./cmd/rpg-service" 9083
start_air "./cmd/feat-service" 9084
start_air "./cmd/update-service" 9085

wait
