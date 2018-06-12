#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.master.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

start() {
  myip=`ifconfig | grep 'inet addr:192.168.212' | awk -F':' '{print $2}' | awk '{print $1}'`
  docker run -d \
    --name alluxio-master \
    --hostname $myip \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_MASTER_HOSTNAME=$myip \
    -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.212.42:2181,192.168.212.45:2181,192.168.212.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -p 19998:19998 \
    -p 19999:19999 \
    -v /alluxio-share/alluxio/journal/:/journal \
    -v /alluxio-share/alluxio/underStorage/:/underStorage \
    alluxio \
    master --no-format
}

remove() {
  docker rm -f alluxio-master
}

status() {
  docker ps -a | grep alluxio-master
  echo "alluxio-master logs:"
  docker logs --tail 25 alluxio-master
}

case $cmd in
  pull)
    docker pull reg-xs.qiniu.io/atlab/alluxio
    docker tag reg-xs.qiniu.io/atlab/alluxio alluxio
  ;;
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
