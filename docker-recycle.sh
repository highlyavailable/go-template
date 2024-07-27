containers=$(docker ps -a --filter "name=goapp" -q)
if [ -n "$containers" ]; then
  docker stop $containers
  docker rm $containers
fi

images=$(docker images --filter "reference=goapp" -q)
if [ -n "$images" ]; then
  docker rmi $images
fi

docker-compose up --build -d