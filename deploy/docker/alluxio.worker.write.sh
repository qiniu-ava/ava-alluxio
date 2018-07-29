#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.worker.writer.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

myip=`ifconfig | grep 'inet addr:192.168.213' | awk -F':' '{print $2}' | awk '{print $1}'`

start() {
  container_mem_size=16g

  mkdir -p /alluxio-share/alluxio/workers/${myip}/tmp
  mkdir -p /alluxio-share/alluxio/workers/${myip}/cachedisk
  mkdir -p /alluxio-share/alluxio/workers/${myip}/underStorage

  source /alluxio-share/alluxio/env/worker

  docker run -d \
    --name alluxio-worker-writer \
    --hostname ${myip} \
    -m ${container_mem_size} \
    -e ALLUXIO_JAVA_OPTS="-Xms8g -Xmx8g -Xss4m" \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_WORKER_PORT=${ALLUXIO_WRITE_WORKER_PORT} \
    -e ALLUXIO_WORKER_DATA_PORT=${ALLUXIO_WRITE_WORKER_DATA_PORT} \
    -e ALLUXIO_WORKER_WEB_PORT=${ALLUXIO_WRITE_WORKER_WEB_PORT} \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=256 \
    -e KODO_IO_ORIGHOST=${KODO_IO_ORIGHOST} \
    -e KODO_UP_ORIGHOST=${KODO_UP_ORIGHOST} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=1 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/cachedisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=4TB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.01 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.002 \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=512 \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=2048 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.212.42:2181,192.168.212.45:2181,192.168.212.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -p ${ALLUXIO_WRITE_WORKER_PORT}:${ALLUXIO_WRITE_WORKER_PORT} \
    -p ${ALLUXIO_WRITE_WORKER_DATA_PORT}:${ALLUXIO_WRITE_WORKER_DATA_PORT} \
    -p ${ALLUXIO_WRITE_WORKER_WEB_PORT}:${ALLUXIO_WRITE_WORKER_WEB_PORT} \
    -v /alluxio-share/alluxio/workers/${myip}/tmp:/tmp \
    -v /alluxio-share/alluxio/workers/${myip}/cachedisk:/opt/cachedisk \
    -v /alluxio-share/alluxio/workers/${myip}/underStorage:/underStorage \
    --restart=always \
    alluxio \
    worker --no-format
}

remove() {
  docker rm -f alluxio-worker-writer
}

status() {
  docker ps -a | grep alluxio-worker-writer
  echo "alluxio-worker-writer logs:"
  docker logs --tail 25 alluxio-worker-writer
}

case $cmd in
  pull)
    tag=$2
    if [ "$tag" = "" ];then
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
