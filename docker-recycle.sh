#!/bin/bash

# Stop and remove all goapp containers
echo "-> Stopping and removing all goapp containers"
containers=$(docker ps -a --filter "name=goapp" -q)
if [ -n "$containers" ]; then
  docker stop $containers
  docker rm $containers
fi

# Remove all goapp images
echo "-> Removing all goapp images"
images=$(docker images --filter "reference=goapp" -q)
if [ -n "$images" ]; then
  docker rmi $images
fi

# Build and run the goapp container
echo "-> Building and running the goapp container"
docker-compose up --build -d