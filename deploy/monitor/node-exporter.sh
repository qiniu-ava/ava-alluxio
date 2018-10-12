#!/bin/bash

cmd=$1
deployname=$2
if [ "$deployname" = "" ]; then
  deployname=node-exporter
fi

if [ "$cmd" = "" ]; then
  echo "usage: ./node-exporter.sh <cmd>, where cmd should be one of start/restart/remove/status"
  exit 1
fi

start() {
  docker run --name=node-exporter \
    -v /:/rootfs:ro \
    --detach=true \
    -p 9100:9100 \
    --restart=always \
    prom/node-exporter
}

remove() {
  docker rm -f node-exporter
}

status() {
  docker ps -a | grep node-exporter
  echo "node-exporter logs:"
  docker logs --tail 25 node-exporter
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
