#!/bin/bash

cmd=$1
deployname=$2
if [ "$deployname" = "" ]; then
  deployname=cadvisor
fi

if [ "$cmd" = "" ]; then
  echo "usage: ./cadvisor.sh <cmd>, where cmd should be one of start/restart/remove/status"
  exit 1
fi

start() {
  docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --publish=8080:8080 \
  --detach=true \
  --name=cadvisor \
  --restart=always \
  google/cadvisor:latest
}

remove() {
  docker rm -f cadvisor
}

status() {
  docker ps -a | grep cadvisor
  echo "cadvisor logs:"
  docker logs --tail 25 cadvisor
}

case $cmd in
  start)
    start
  ;;
  restart)
    remove
    start
  ;;
  remove)
    remove
  ;;
  rm)
    remove
  ;;
  status)
    status
  ;;
esac