#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/taskSvc ./cmd/task-service 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/rpgSvc ./cmd/rpg-service 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/updateSvc ./cmd/update-service 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/featSvc ./cmd/feat-service 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/userSvc ./cmd/user-service 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/authSvc ./cmd/auth-service 
