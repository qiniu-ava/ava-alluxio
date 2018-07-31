#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ ! -f "$DIR/env" ]]; then
  echo "plz config your local volume paths in path $DIR/env"
  exit 1
fi

# shellcheck disable=SC1090

source "$DIR"/env
myip=$(ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | awk '{print $2}')
docker rm -f alluxio-master
docker run -d \
  --name alluxio-master \
  --hostname "$myip" \
  -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
  -e ALLUXIO_MASTER_HOSTNAME="$myip" \
  -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
  -e ALLUXIO_MASTER_WORKER_TIMEOUT=24h \
  -e ALLUXIO_MASTER_PORT="${ALLUXIO_MASTER_PORT}" \
  -e ALLUXIO_MASTER_WEB_PORT="${ALLUXIO_MASTER_WEB_PORT}" \
  -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
  -p "${ALLUXIO_MASTER_PORT}":"${ALLUXIO_MASTER_PORT}" \
  -p "${ALLUXIO_MASTER_WEB_PORT}":"${ALLUXIO_MASTER_WEB_PORT}" \
  -v "${ALLUXIO_MASTER_JOURNAL}":/journal \
  -v "${ALLUXIO_MASTER_UNDERFSSTORAGE}":/underStorage \
  --restart=always \
  alluxio \
  master --no-format
