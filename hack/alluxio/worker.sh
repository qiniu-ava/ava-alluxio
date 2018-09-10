#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ ! -f "$DIR/env" ]]; then
  echo "plz config your setting in path $DIR/env, env.template may be helpful for you"
  exit 1
fi

# shellcheck disable=SC1090
source "$DIR"/env
myip=$(ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | awk '{print $2}')
docker rm -f alluxio-worker
docker run -d \
  --name alluxio-worker \
  --hostname "${myip}" \
  -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
  -e ALLUXIO_RAM_FOLDER=/opt/ramdisk \
  -e KODO_IO_ORIGHOST="${KODO_IO_ORIGHOST}" \
  -e KODO_UP_ORIGHOST="${KODO_UP_ORIGHOST}" \
  -e ALLUXIO_WORKER_MEMORY_SIZE="${ALLUXIO_WORKER_MEMORY_SIZE}" \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVELS=2 \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_ALIAS=MEM \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_PATH=/opt/ramdisk \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_DIRS_QUOTA="${ALLUXIO_WORKER_MEMORY_SIZE}" \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_HIGH_RATIO=0.5 \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL0_WATERMARK_LOW_RATIO=0.3 \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_ALIAS=SSD \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_PATH=/opt/cachedisk \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_DIRS_QUOTA=10GB \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_HIGH_RATIO=0.8 \
  -e ALLUXIO_WORKER_TIEREDSTORE_LEVEL1_WATERMARK_LOW_RATIO=0.7 \
  -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_ENABLED=true \
  -e ALLUXIO_WORKER_TIEREDSTORE_RESERVER_INTERVAL=10000 \
  -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
  -e ALLUXIO_WORKER_PORT="${ALLUXIO_WORKER_PORT}" \
  -e ALLUXIO_WORKER_DATA_PORT="${ALLUXIO_WORKER_DATA_PORT}" \
  -e ALLUXIO_WORKER_WEB_PORT="${ALLUXIO_WORKER_WEB_PORT}" \
  -e ALLUXIO_MASTER_HOSTNAME="${myip}" \
  -p "${ALLUXIO_WORKER_PORT}":"${ALLUXIO_WORKER_PORT}" \
  -p "${ALLUXIO_WORKER_DATA_PORT}":"${ALLUXIO_WORKER_DATA_PORT}" \
  -p "${ALLUXIO_WORKER_WEB_PORT}":"${ALLUXIO_WORKER_WEB_PORT}" \
  -v "${ALLUXIO_WORKER_RAM_DISK}":/opt/ramdisk \
  -v "${ALLUXIO_WORKER_SSD_DISK}":/opt/cachedisk \
  --restart=always \
  alluxio \
  worker --no-format
