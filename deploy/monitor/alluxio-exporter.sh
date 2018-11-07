#!/bin/bash

cmd=$1
group=$2
master_host=$3
if [ "$cmd" = "" ] || [ "$group" == "" ] || [ "$master_host" == "" ] ; then
  echo "usage: ./exporter-start.sh <cmd> <group_name> <cluster_master_ingress or master_host:webport> where cmd should be one of start/restart/remove/status"
  exit 1
fi

start() {
  if [ "$group" == "" ] || [ "$master_host" == "" ] ; then
    echo "usage: ./exporter-start.sh start <group_name> <cluster_master_ingress or master_host:webport>"
    return
  fi
  python exporter_conf.py "$group" "$master_host"
  docker run \
    --publish=9999:9999 \
    --detach=true \
    --name=alluxio-exporter-"$group" \
    --restart=always \
    -v "$HOME"/alluxio-exporter/:/conf \
    reg-xs.qiniu.io/atlab/alluxio-exporter:v0.1.0 \
    --exporter.config=/conf/"$group"-exporter.yml
}

remove() {
  if [ "$group" == "" ] ; then
    echo "usage: ./exporter-start.sh remove <group_name>"
    return
  fi
  docker rm -f alluxio-exporter-"$group"
}

status() {
  docker ps -a | grep alluxio-exporter-"$group"
  echo "cadvisor logs:"
  docker logs --tail 25 alluxio-exporter-"$group"
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
