#!/bin/bash

BYellow='\033[1;33m'
BRed='\033[1;31m'
NC="\033[0m"

echo -e "${BRed}[ERROR]${NC} this script is ${BRed}Depricated${NC}, please try to deploy zookeeper by yourself"
exit 1

set -ex

source /root/env

if [ "$ZOO_SERVERS" = "" ]; then
  echo "ZOO_SERVERS not set"
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

myip=`ifconfig | grep 'inet addr:192.168.213' | awk -F':' '{print $2}' | awk '{print $1}'`
zoo_servers=`echo "$ZOO_SERVERS" | sed "s|\$myip|0.0.0.0|"`

mkdir -m 777 -p /disk1/zk/$myip/config /disk1/zk/$myip/data /disk1/zk/$myip/log
rm -rf /disk1/zk/$myip/data/*
rm -rf /disk1/zk/$myip/config/zoo.cfg
cp ${DIR}/zookeeper/configuration.xsl /disk1/zk/$myip/config
cp ${DIR}/zookeeper/log4j.properties /disk1/zk/$myip/config

docker rm -f alluxio-zk
docker run -d \
  --network host \
  --name alluxio-zk \
  -e ZOO_MY_ID=${ZK_MY_ID} \
  -e ZOO_MAX_CLIENT_CNXNS=3600 \
  -e ZOO_SERVERS=${zoo_servers} \
  -v /disk1/zk/$myip/config:/conf \
  -v /disk1/zk/$myip/data:/data \
  -v /disk1/zk/$myip/log:/datalog \
  -p 2181:2181 \
  -p 2888:2888 \
  -p 3888:3888 \
  --restart=always \
  zookeeper:3.4
