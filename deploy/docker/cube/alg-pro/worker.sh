#!/bin/bash

######################################################################
# worker node list:
# jq19
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. ${DIR}/../common/util.sh

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.worker.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

myip=$(getMyIP)
group=alg-pro

start() {
  ram_tier_size=10G

  # implements in common/util.sh
  ensureGroupRamdisk $group $ram_tier_size

  mkdir -p /disk1/alluxio/data-$group/cachedisk
  mkdir -p /disk2/alluxio/data-$group/cachedisk
  mkdir -p /disk3/alluxio/data-$group/cachedisk
  mkdir -p /disk4/alluxio/data-$group/cachedisk
  mkdir -p /disk5/alluxio/data-$group/cachedisk

  source /disk-cephfs/alluxio/env/worker-$group

  docker run -d \
    --name alluxio-worker-$group \
    --hostname ${myip} \
    --network host \
    -e ALLUXIO_JAVA_OPTS="-Xmx8g -XX:+UseG1GC " \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=256 \
    -e KODO_IO_ORIGHOST=${KODO_IO_ORIGHOST} \
    -e KODO_UP_ORIGHOST=${KODO_UP_ORIGHOST} \
    -e ALLUXIO_WORKER_PORT=${ALLUXIO_WORKER_PORT} \
    -e ALLUXIO_WORKER_DATA_PORT=${ALLUXIO_WORKER_DATA_PORT} \
    -e ALLUXIO_WORKER_WEB_PORT=${ALLUXIO_WORKER_WEB_PORT} \
    -e ALLUXIO_WORKER_MEMORY_SIZE=${ram_tier_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=2 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=${ram_tier_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.7 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.5 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk1,/opt/cachedisk2,/opt/cachedisk3,/opt/cachedisk4,/opt/cachedisk5 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=200GB,200GB,200GB,200GB,200GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=0.8 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=0.7 \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=512 \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=4096 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/$group \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/$group \
    -v /mnt/ramdisk-$group:/opt/ramdisk \
    -v /disk1/alluxio/data-$group/cachedisk:/opt/cachedisk1 \
    -v /disk2/alluxio/data-$group/cachedisk:/opt/cachedisk2 \
    -v /disk3/alluxio/data-$group/cachedisk:/opt/cachedisk3 \
    -v /disk4/alluxio/data-$group/cachedisk:/opt/cachedisk4 \
    -v /disk5/alluxio/data-$group/cachedisk:/opt/cachedisk5 \
    --restart=always \
    alluxio-${group} \
    worker --no-format
}

remove() {
  docker rm -f alluxio-worker-$group
}

status() {
  docker ps -a | grep alluxio-worker-$group
  echo "alluxio-worker logs:"
  docker logs --tail 25 alluxio-worker-$group
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
    docker tag reg-xs.qiniu.io/atlab/alluxio:$tag alluxio-${group}
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
