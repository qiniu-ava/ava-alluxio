#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./cadvisor.sh <cmd>, where cmd should be one of start/restart/remove/status"
  exit 1
fi

start() {
  myip=$(ifconfig | grep '100.100.35' | awk '{print $2}')
  docker run \
    --publish=9998:9998 \
    --detach=true \
    --name=jvm-exporter \
    --restart=always \
    -e HOSTIP="${myip}" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    reg-xs.qiniu.io/atlab/jvm-exporter:v0.1.0
}

remove() {
  docker rm -f jvm-exporter
}

status() {
  docker ps -a | grep jvm-exporter
  echo "cadvisor logs:"
  docker logs --tail 25 kvm-exporter
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
