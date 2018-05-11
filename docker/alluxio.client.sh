#!/bin/bash

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.client.sh <cmd>, where cmd should be one of pull/start/restart/remove/status"
  exit 1
fi

start() {
  docker run -d \
    --name alluxio-client \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
    -e ALLUXIO_MASTER_PORT=19998 \
    -e ALLUXIO_WORKER_BLOCK_MASTER_CLIENT_POOL_SIZE=64 \
    -e ALLUXIO_WORKER_MEMORY_SIZE=10GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=2 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA=10GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.75 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.5 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=50GB \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=0.75 \
    -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=0.5 \
    -e ALLUXIO_FUSE_CACHED_PATHS_MAX=5000 \
    -e ALLUXIO_USER_BLOCK_MASTER_CLIENT_THREADS=256 \
    -e ALLUXIO_USER_BLOCK_WORKER_CLIENT_THREADS=256 \
    -e ALLUXIO_USER_FILE_MASTER_CLIENT_THREADS=156 \
    -e ALLUXIO_USER_BLOCK_SIZE_BYTES_DEFAULT=1MB \
    -e ALLUXIO_USER_NETWORK_NETTY_WORKER_THREADS=8192 \
    -e ALLUXIO_SECURITY_GROUP_MAPPING_CLASS="" \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/o    khttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/    jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=10.200.20.91:2181,10.200.20.70:2181,10.200.20.80:2181 \
    -p 29998:29998 \
    -p 39998:39998 \
    alluxio
}

remove() {
  docker rm -f alluxio-client
}

status() {
  docker ps | grep alluxio-client
  echo "alluxio-client logs:"
  docker logs --tail 25 alluxio-client
}

case $cmd in
  pull)
    docker pull reg-xs.qiniu.io/atlab/alluxio-bowen
    docker tag reg-xs.qiniu.io/atlab/alluxio-bowen alluxio
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
