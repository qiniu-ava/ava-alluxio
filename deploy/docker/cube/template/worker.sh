#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. ${DIR}/../common/util.sh

# GROUP should be set
if [ "${GROUP}" == "" ]; then
  echo -e "${BRed}[FATAL]${NC} GROUP not set"
  exit 1
fi

# check if this node is valid for <GROUP>-worker
if ! echo "$NODE_LIST" | grep -wE $(hostname) > /dev/null ; then
  echo -e "${BRed}[FATAL]${NC} $(hostname) is not in worker node list" 
  exit 1
fi

cmd=$1
container_name="alluxio-worker-${GROUP}"

if [ "$cmd" = "" ]; then
  echo -e "${BRed}[FATAL]${NC} usage: ./alluxio.worker.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

myip=$(getMyIP)

jvm_size=8G
ram_tier_size=40G
ram_tier_high=0.7
ram_tier_low=0.5
ssd_tier_high=0.8
ssd_tier_low=0.7

if [ "${WORKER_JVM_SIZE}" != "" ]; then
  jvm_size=${WORKER_JVM_SIZE}
fi

if [ "${WORKER_RAM_TIER_SIZE}" != "" ]; then
  ram_tier_size=${WORKER_RAM_TIER_SIZE}
fi

if [ "${WORKER_RAM_TIER_HIGHT}" != "" ]; then
  ram_tier_high=${WORKER_RAM_TIER_HIGHT}
fi

if [ "${WORKER_RAM_TIER_LOW}" != "" ]; then
  ram_tier_low=${WORKER_RAM_TIER_LOW}
fi

if [ "${WORKER_SSD_TIER_HIGHT}" != "" ]; then
  ssd_tier_high=${WORKER_SSD_TIER_HIGHT}
fi

if [ "${WORKER_SSD_TIER_LOW}" != "" ]; then
  ssd_tier_high=${WORKER_SSD_TIER_LOW}
fi

start() {
  # implements in common/util.sh
  ensureGroupRamdisk ${GROUP} $ram_tier_size

  ssd=$(getAvailableSSD)

  if [ "${#ssd[@]}" -eq 0 ]; then
    echo -e "${BRed}[Fatal]${NC} no available ssd for worker"
    exit 1
  fi

  for disk in "${ssd}"; do
    mkdir -p "${disk}/alluxio/data-${GROUP}/cachedisk"
  done

  # -v /disk1/alluxio/data-${GROUP}/cachedisk:/opt/cachedisk1 \
  volume_str=$(gen_volume_str_from_ssd_list ${GROUP} ${ssd})
  # /opt/cachedisk1,/opt/cachedisk2,/opt/cachedisk3,/opt/cachedisk4,/opt/cachedisk5
  path_str=$(gen_path_str_from_ssd_list ${ssd})
  # 200GB,200GB,200GB,200GB,200GB
  quota_str=$(gen_quota_str_from_ssd_list ${ssd})

  . /disk-cephfs/alluxio/env/worker-${GROUP}

  docker run -d \
    --name alluxio-worker-${GROUP} \
    --hostname ${myip} \
    --network host \
    -e ALLUXIO_JAVA_OPTS="-Xmx${jvm_size} -XX:+UseG1GC " \
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
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=${ram_tier_high} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=${ram_tier_low} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH="${path_str}" \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA="${quota_str}" \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=${ssd_tier_high} \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=${ssd_tier_low} \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
    -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/${GROUP} \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/${GROUP} \
    -v /mnt/ramdisk-${GROUP}:/opt/ramdisk \
    ${volume_str} \
    --restart=always \
    alluxio-${GROUP} \
    worker --no-format
}

remove() {
  docker rm -f alluxio-worker-${GROUP}
}

status() {
  docker ps -a | grep alluxio-worker-${GROUP}
  echo "alluxio-worker logs:"
  docker logs --tail 25 alluxio-worker-${GROUP}
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
