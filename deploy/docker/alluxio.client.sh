#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.client.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

start() {
  docker run -d \
    --name alluxio-client \
    --privileged=true \
    -e ALLUXIO_FUSE_CACHED_PATHS_MAX=5000 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=256 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=256 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=156 \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=8192 \
    -e ALLUXIO_USER_FILE_WRITE_TIER_DEFAULT=1 \
    -e ALLUXIO_USER_FILE_WRITETYPE_DEFAULT=ASYNC_THROUGH \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -e QINIU_WRITER_HOSTS="192.168.213.50:29899,192.168.213.51:29899,192.168.213.52:29899" \
    -e QINIU_READER_EX_HOSTS="192.168.213.50:29899,192.168.213.51:29899,192.168.213.52:29899" \
    alluxio \
    proxy
}

remove() {
  docker rm -f alluxio-client
}

status() {
  docker ps -a | grep alluxio-client
  echo "alluxio-client logs:"
  docker logs --tail 25 alluxio-client
}

case $cmd in
  pull)
    tag=$2
    if [ $tag = "" ];then
      cd $DIR/../../alluxio && alluxio_hash=`git rev-parse --short=7 HEAD` && cd -
      cd $DIR/../../kodo && kodo_hash=`git rev-parse --short=7 HEAD` && cd -
      tag=$alluxio_hash-$kodo_hash
    fi
    docker pull reg-xs.qiniu.io/atlab/alluxio:$tag
    docker tag reg-xs.qiniu.io/atlab/alluxio:$tag alluxio
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
