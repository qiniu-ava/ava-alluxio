#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.worker.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

myip=`ifconfig | grep 'inet addr:192.168.213' | awk -F':' '{print $2}' | awk '{print $1}'`

start() {
  ram_size=180G
  ram_tier_size=170G
  if [ ! -d /mnt/ramdisk ]; then
    sudo mkdir -p /mnt/ramdisk
    sudo mount -t ramfs -o size=${ram_size} ramfs /mnt/ramdisk
    sudo chmod a+w /mnt/ramdisk
    mkdir -p /mnt/ramdisk/data
  fi

  mkdir -p /disk1/alluxio/data/cachedisk
  mkdir -p /disk2/alluxio/data/cachedisk
  mkdir -p /disk2/alluxio/data/underStorage

  docker run -d \
    --name alluxio-worker \
    --hostname ${myip} \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=256 \
    -e KODO_ORIGHOST=http://nbjjh-gate-io.qiniu.com \
    -e ALLUXIO_WORKER_MEMORY_SIZE=$ram_size \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=3 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=$ram_tier_size \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.75 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.5 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk1 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=300GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=0.9 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=0.7 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL2_ALIAS=HDD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL2_DIRS_PATH=/opt/cachedisk2 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL2_DIRS_QUOTA=400GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL2_WATERMARK_HIGH_RATIO=0.9 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL2_WATERMARK_LOW_RATIO=0.7 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=2048 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=2048 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=8192 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.212.42:2181,192.168.212.45:2181,192.168.212.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -p 29998:29998 \
    -p 29999:29999 \
    -p 30000:30000 \
    -p 39998:39998 \
    -v /mnt/ramdisk:/opt/ramdisk \
    -v /disk1/alluxio/data/cachedisk:/opt/cachedisk1 \
    -v /disk2/alluxio/data/cachedisk:/opt/cachedisk2 \
    -v /alluxio-share/alluxio/underStorage:/underStorage \
    alluxio \
    worker
}

remove() {
  docker rm -f alluxio-worker
}

status() {
  docker ps -a | grep alluxio-worker
  echo "alluxio-worker logs:"
  docker logs --tail 25 alluxio-worker
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
