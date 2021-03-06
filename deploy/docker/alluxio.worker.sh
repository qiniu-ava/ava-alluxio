#!/bin/bash

BYellow='\033[1;33m'
BRed='\033[1;31m'
NC="\033[0m"

echo -e "${BYellow}[WARNING]${NC} this script is ${BRed}Depricated${NC}, please refer to README"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.worker.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

myip=`ifconfig | grep 'inet addr:192.168.213' | awk -F':' '{print $2}' | awk '{print $1}'`

start() {
  container_mem_size=130g # worker 容器运行时申请的内存上线
  worker_mem_size=100G
  ram_tier_size=100G
  if [ ! -d /mnt/ramdisk ]; then
    sudo mkdir -p /mnt/ramdisk
    sudo mount -t ramfs -o size=${ram_tier_size} ramfs /mnt/ramdisk
    sudo chmod a+w /mnt/ramdisk
    mkdir -p /mnt/ramdisk/data
  fi

  mkdir -p /disk9/alluxio/data/cachedisk
  mkdir -p /disk10/alluxio/data/cachedisk
  mkdir -p /disk11/alluxio/data/cachedisk
  mkdir -p /disk12/alluxio/data/cachedisk

  source /alluxio-share/alluxio/env/worker

  docker run -d \
    --name alluxio-worker \
    --hostname ${myip} \
    --network host \
    -m ${container_mem_size} \
    -e ALLUXIO_JAVA_OPTS="-Xms16g" \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=256 \
    -e KODO_IO_ORIGHOST=${KODO_IO_ORIGHOST} \
    -e KODO_UP_ORIGHOST=${KODO_UP_ORIGHOST} \
    -e ALLUXIO_WORKER_MEMORY_SIZE=${worker_mem_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=2 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=${ram_tier_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.7 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.5 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk1,/opt/cachedisk2,/opt/cachedisk3,/opt/cachedisk4 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=300GB,400GB,1800GB,1800GB,1800GB,1800GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=0.8 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=0.7 \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=2048 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=2048 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=1024 \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=8192 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -v /mnt/ramdisk:/opt/ramdisk \
    -v /disk9/alluxio/data/cachedisk:/opt/cachedisk1 \
    -v /disk10/alluxio/data/cachedisk:/opt/cachedisk2 \
    -v /disk11/alluxio/data/cachedisk:/opt/cachedisk3 \
    -v /disk12/alluxio/data/cachedisk:/opt/cachedisk4 \
    -v /alluxio-share/alluxio/underStorage:/underStorage \
    --restart=always \
    alluxio \
    worker --no-format
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
