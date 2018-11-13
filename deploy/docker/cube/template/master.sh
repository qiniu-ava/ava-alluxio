#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. ${DIR}/../common/util.sh

# GROUP should be set
if [ "${GROUP}" == "" ]; then
  echo -e "${BRed}[FATAL]${NC} GROUP not set"
  exit 1
fi

# check if this node is valid for <GROUP>-master
if ! echo "$NODE_LIST" | grep -wE $(hostname) > /dev/null ; then
  echo -e "${BRed}[FATAL]${NC} $(hostname) is not in master node list" 
  exit 1
fi

cmd=$1
container_name="alluxio-master-${GROUP}"

if [ "$cmd" = "" ]; then
  echo "usage: ./alluxio.master.sh <cmd> [options]"
  echo "  where cmd should be one of pull/start/restart/remove/status"
  echo "  options:"
  echo "    pull [tag] default tag will be <hashofalluxio-hashofkodo>"
  exit 1
fi

. /disk-cephfs/alluxio/env/master-"${GROUP}"

jvm_size=48G
inode_capacity=3000000
inode_evict_ratio=80
startup_check_consistency=false

if [ "${MASTER_JVM_SIZE}" != "" ]; then
  jvm_size=${MASTER_JVM_SIZE}
fi

if [ "${MASTER_INODE_CAPACITY}" != "" ]; then
  inode_capacity=${MASTER_INODE_CAPACITY}
fi

if [ "${MASTER_INODE_EVICT_RATIO}" != "" ]; then
  inode_evict_ratio=${MASTER_INODE_EVICT_RATIO}
fi

if [ "${MASTER_STARTUP_CONSISTENCY_CHECK}" != "" ]; then
  startup_check_consistency=${MASTER_STARTUP_CONSISTENCY_CHECK}
fi

start() {
  myip=$(getMyIP)
  docker run -d \
    --name ${container_name} \
    --hostname $myip \
    --network host \
    -e ALLUXIO_JAVA_OPTS=" -Xmx${jvm_size} -XX:+UseG1GC " \
    -e ALLUXIO_UNDERFS_ADDRESS=/underStorage \
    -e ALLUXIO_MASTER_HOSTNAME=$myip \
    -e ALLUXIO_MASTER_UFS_PATH_CACHE_THREADS=0 \
    -e ALLUXIO_MASTER_JOURNAL_FOLDER=/journal \
    -e ALLUXIO_MASTER_STARTUP_CONSISTENCY_CHECK_ENABLED=${startup_check_consistency} \
    -e ALLUXIO_MASTER_JOURNAL_CHECKPOINT_PERIOD_ENTRIES=500000 \
    -e ALLUXIO_MASTER_STARTUP_BLOCK_INTEGRITY_CHECK_ENABLED=true \
    -e ALLUXIO_MASTER_INODE_CAPACITY=${inode_capacity} \
    -e ALLUXIO_MASTER_INODE_EVICT_RATIO=${inode_evict_ratio} \
    -e ALLUXIO_CLASSPATH=/opt/alluxio/lib/gson-2.2.4.jar:/opt/alluxio/lib/qiniu-java-sdk-7.2.11.jar:/opt/alluxio/lib/okhttp-3.10.0.jar:/opt/alluxio/lib/okio-1.14.0.jar:/opt/alluxio/lib/jackson-databind-2.9.5.jar:/opt/alluxio/lib/jackson-core-2.9.5.jar:/opt/alluxio/lib/jackson-annotations-2.9.5.jar \
    -e ALLUXIO_ZOOKEEPER_ENABLED=true \
    -e ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181 \
    -e ALLUXIO_ZOOKEEPER_LEADER_PATH=/leader/$GROUP \
    -e ALLUXIO_ZOOKEEPER_ELECTION_PATH=/election/$GROUP \
    -v /disk-cephfs/alluxio/$GROUP/journal/:/journal \
    -v /disk-cephfs/alluxio/$GROUP/underStorage/:/underStorage \
    --restart=always \
    alluxio-${GROUP} \
    master --no-format
}

remove() {
  docker rm -f ${container_name}
}

status() {
  docker ps -a | grep ${container_name}
  echo "alluxio-master logs:"
  docker logs --tail 25 ${container_name}
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
  status)
    status
  ;;
esac
