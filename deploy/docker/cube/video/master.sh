#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1
group=video

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.master.sh <cmd> [options]"
  echo "  where cmd should be one of pull/start/restart/remove/status"
  echo "  options:"
  echo "    pull [tag] default tag will be <hashofalluxio-hashofkodo>"
  exit 1
fi

mkdir -p /disk-cephfs/alluxio/$group/journal /disk-cephfs/alluxio/$group/underStorage

start() {
  myip=`ifconfig | grep 'inet addr:192.168.213' | awk -F':' '{print $2}' | awk '{print $1}'`
  docker run -d \
    --name alluxio-master-$group \
    --hostname $myip \
    --network host \
    -e ALLUXIO_JAVA_OPTS=" -Xmx72g -XX:+UseG1GC " \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_MASTER_HOSTNAME=$myip \
    -e ALLUXIO_MASTER_UFS_PATH_CACHE_THREADS=0 \
    -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
    -e ALLUXIO_MASTER_STARTUP_CONSISTENCY_CHECK_ENABLED=false \
    -e ALLUXIO_MASTER_JOURNAL_CHECKPOINT_PERIOD_ENTRIES=500000 \
    -e ALLUXIO_MASTER_STARTUP_BLOCK_INTEGRITY_CHECK_ENABLED=true \
    -e ALLUXIO_MASTER_INODE_CAPACITY=5000000 \
    -e ALLUXIO_MASTER_INODE_EVICT_RATIO=85 \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/$group \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/$group \
    -v /disk-cephfs/alluxio/$group/journal/:/journal \
    -v /disk-cephfs/alluxio/$group/underStorage/:/underStorage \
    --restart=always \
    alluxio \
    master --no-format
}

remove() {
  docker rm -f alluxio-master-$group
}

status() {
  docker ps -a | grep alluxio-master-$group
  echo "alluxio-master logs:"
  docker logs --tail 25 alluxio-master-$group
}

case $cmd in
  pull)
    tag=$2
    if [ "$tag" = "" ];then
      cd $DIR/../../alluxio && alluxio_hash=`git rev-parse --short=7 HEAD` && cd -
      cd $DIR/../../kodo && kodo_hash=`git rev-parse --short=7 HEAD` && cd -
      tag=$alluxio_hash-$kodo_hash
    fi
    tag="2bece64-080ef7c"
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
  status)
    status
  ;;
esac
