#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./cadvisor.sh <cmd>, where cmd should be one of start/restart/remove/status"
  exit 1
fi

start() {
  docker run \
    --publish=9999:9999 \
    --detach=true \
    --name=alluxio-exporter \
    --restart=always \
    -v "$HOME"/alluxio-exporter/:/conf \
    alluxio-exporter:test \
    --exporter.config=/conf/exporter.yml
}

remove() {
  docker rm -f alluxio-exporter
}

status() {
  docker ps -a | grep alluxio-exporter
  echo "cadvisor logs:"
  docker logs --tail 25 alluxio-exporter
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
