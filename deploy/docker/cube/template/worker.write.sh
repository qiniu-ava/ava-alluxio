#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. ${DIR}/../common/util.sh

# GROUP should be set
if [ "${GROUP}" == "" ]; then
  echo -e "${BRed}[FATAL]${NC} GROUP not set"
  exit 1
fi

# check if this node is valid for <GROUP>-worker-write
if ! echo "$NODE_LIST" | grep -wE $(hostname) > /dev/null ; then
  echo -e "${BRed}[FATAL]${NC} $(hostname) is not in worker-write node list" 
  exit 1
fi

cmd=$1
container_name="alluxio-worker-write-${GROUP}"

if [ "$cmd" = "" ]; then
  echo -e "${BRed}[FATAL]${NC} usage: ./alluxio.worker.write.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

. /disk-cephfs/alluxio/env/worker-${GROUP}

myip=$(getMyIP)

jvm_size=8G
ram_tier_size=4G
ram_tier_high=0.2
ram_tier_low=0.1
ssd_tier_quota=500GB
ssd_tier_high=0.1
ssd_tier_low=0.05

if [ "${WRITE_WORKER_JVM_SIZE}" != "" ]; then
  jvm_size=${WRITE_WORKER_JVM_SIZE}
fi

if [ "${WRITE_WORKER_RAM_TIER_SIZE}" != "" ]; then
  ram_tier_size=${WRITE_WORKER_RAM_TIER_SIZE}
fi

if [ "${WRITE_WORKER_RAM_TIER_HIGHT}" != "" ]; then
  ram_tier_high=${WRITE_WORKER_RAM_TIER_HIGHT}
fi

if [ "${WRITE_WORKER_RAM_TIER_LOW}" != "" ]; then
  ram_tier_low=${WRITE_WORKER_RAM_TIER_LOW}
fi

if [ "${WRITE_WORKER_SSD_TIER_QUOTA}" != "" ]; then
  ssd_tier_quota=${WRITE_WORKER_SSD_TIER_QUOTA}
fi

if [ "${WRITE_WORKER_SSD_TIER_HIGHT}" != "" ]; then
  ssd_tier_high=${WRITE_WORKER_SSD_TIER_HIGHT}
fi

if [ "${WRITE_WORKER_SSD_TIER_LOW}" != "" ]; then
  ssd_tier_high=${WRITE_WORKER_SSD_TIER_LOW}
fi

start() {
  ensureGroupRamdisk "${GROUP}" $ram_tier_size write

  mkdir -p /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/tmp
  mkdir -p /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/cachedisk
  mkdir -p /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/underStorage

  docker run -d \
    --name $container_name \
    --hostname ${myip} \
    --network host \
    -e ALLUXIO_JAVA_OPTS="-Xmx${jvm_size} -XX:+UseG1GC " \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
    -e ALLUXIO_WORKER_PORT=${ALLUXIO_WRITE_WORKER_PORT} \
    -e ALLUXIO_WORKER_DATA_PORT=${ALLUXIO_WRITE_WORKER_DATA_PORT} \
    -e ALLUXIO_WORKER_WEB_PORT=${ALLUXIO_WRITE_WORKER_WEB_PORT} \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=256 \
    -e KODO_IO_ORIGHOST=${KODO_IO_ORIGHOST} \
    -e KODO_UP_ORIGHOST=${KODO_UP_ORIGHOST} \
    -e ALLUXIO_WORKER_MEMORY_SIZE=${ram_tier_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=2 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=${ram_tier_size} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=${ram_tier_high} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=${ram_tier_low} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=${ssd_tier_quota} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=${ssd_tier_high} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=${ssd_tier_low} \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/"${GROUP}" \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/"${GROUP}" \
    -v /mnt/ramdisk-write-"${GROUP}":/opt/ramdisk \
    -v /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/tmp:/tmp \
    -v /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/cachedisk:/opt/cachedisk \
    -v /disk-rbd/alluxio/workers/worker-write-"${GROUP}"/${myip}/underStorage:/underStorage \
    --restart=always \
    alluxio-${GROUP} \
    worker --no-format
}

remove() {
  docker rm -f $container_name
}

status() {
  docker ps -a | grep $container_name
  echo "$container_name logs:"
  docker logs --tail 25 $container_name
}

case $cmd in
  pull)
    tag=$2
    if [ "$tag" = "" ];then
      cd $DIR/../../../../alluxio && alluxio_hash=`git rev-parse --short=7 HEAD` && cd -
      cd $DIR/../../../../kodo && kodo_hash=`git rev-parse --short=7 HEAD` && cd -
      tag=$alluxio_hash-$kodo_hash
    fi
    docker pull reg-xs.qiniu.io/atlab/alluxio:$tag
    docker tag reg-xs.qiniu.io/atlab/alluxio:$tag alluxio-${GROUP}
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
