#!/bin/sh

IMAGE_NAME="go_vision"
CONTAINER_NAME="go-vision"
PORT_NUMBER="4000"

export GOOGLE_APPLICATION_CREDENTIALS="/home/username/service_credentials.json"

docker build -t $IMAGE_NAME .

docker run -e ENV=$ENV -d -p $PORT_NUMBER:$PORT_NUMBER -it $IMAGE_NAME

