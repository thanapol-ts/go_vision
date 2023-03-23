#!/bin/sh

docker ps -al

docker stop go_vision

docker rm go_vision

docker rmi go_vision-go:latest

docker compose build

docker compose up -d
