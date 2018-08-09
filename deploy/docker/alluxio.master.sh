#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.master.sh <cmd> [options]"
  echo "  where cmd should be one of pull/start/restart/remove/status"
  echo "  options:"
  echo "    pull [tag] default tag will be <hashofalluxio-hashofkodo>"
  exit 1
fi

start() {
  # master 容器运行时申请的内存上限，分配给 jvm 的内存上限为 64g，
  # 另加上最多 2048 个线程，每个线程 4m 的栈，master 进程最多可占
  # 用 72g 内存
  container_mem_size=75g
  myip=`ifconfig | grep 'inet addr:192.168.212' | awk -F':' '{print $2}' | awk '{print $1}'`
  docker run -d \
    --name alluxio-master \
    --hostname $myip \
    --network host \
    -m ${container_mem_size} \
    -e ALLUXIO_JAVA_OPTS="-Xms64g -Xmx64g -Xss4m" \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_MASTER_HOSTNAME=$myip \
    -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
    -e ALLUXIO_MASTER_WORKER_TIMEOUT=15min \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.212.42:2181,192.168.212.45:2181,192.168.212.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/alluxio-ro \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/alluxio-ro \
    -v /alluxio-share/alluxio/journal/:/journal \
    -v /alluxio-share/alluxio/underStorage/:/underStorage \
    --restart=always \
    alluxio \
    master --no-format
}

remove() {
  docker rm -f alluxio-master
}

status() {
  docker ps -a | grep alluxio-master
  echo "alluxio-master logs:"
  docker logs --tail 25 alluxio-master
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
