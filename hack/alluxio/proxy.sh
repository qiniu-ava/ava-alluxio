#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ ! -f "$DIR/env" ]]; then
  echo "plz config your setting in path $DIR/env, env.template may be helpful for you"
  exit 1
fi

# shellcheck disable=SC1090
source "$DIR"/env
myip=$(ifconfig | grep -E -A6 "^en0:|^bond0" | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1)
docker rm -f alluxio-proxy
docker run -d \
  --name alluxio-proxy \
  --hostname "$myip" \
  --privileged=true \
  -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
  -e ALLUXIO_MASTER_HOSTNAME="$myip" \
  -e ALLUXIO_PROXY_WEB_PORT="${ALLUXIO_PROXY_WEB_PORT}" \
  -p "${ALLUXIO_PROXY_WEB_PORT}":"${ALLUXIO_PROXY_WEB_PORT}" \
  -v "${ALLUXIO_WORKER_SSD_DISK}":/opt/cachedisk-read \
  -v "${ALLUXIO_WRITE_WORKER_SSD_DISK}":/opt/cachedisk-write \
  --restart=always \
  alluxio \
  proxy
