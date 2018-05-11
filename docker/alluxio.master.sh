#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.master.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

start() {
  docker run -d \
    --name alluxio-master \
    -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=10.200.20.91:2181,10.200.20.70:2181,10.200.20.80:2181 \
    -p 19998:19998 \
    -v /alluxio-journal/volumes/test/journal/:/journal \
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
    docker pull reg-xs.qiniu.io/atlab/alluxio-bowen
    docker tag reg-xs.qiniu.io/atlab/alluxio-bowen alluxio
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
