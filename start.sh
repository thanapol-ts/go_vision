#!/bin/sh

docker compose build

docker compose up -d

docker ps -al

