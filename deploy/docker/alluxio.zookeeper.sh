#!/bin/bash

set -ex

source /root/env

if [ $ZOO_SERVERS = "" ]; then
  echo "ZOO_SERVERS not set"
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

myip=`ifconfig | grep 'inet addr:192.168.212' | awk -F':' '{print $2}' | awk '{print $1}'`
myid=`echo $myip | awk -F'.' '{print $4}'`
zoo_servers=`echo $ZOO_SERVERS | sed "s|\$myip|0.0.0.0|"`

mkdir -m 777 -p /alluxio-share/zookeeper/$myip/config /alluxio-share/zookeeper/$myip/data /alluxio-share/zookeeper/$myip/log
rm -rf /alluxio-share/zookeeper/$myip/data/*
rm -rf /alluxio-share/zookeeper/$myip/config/zoo.cfg
cp ${DIR}/zookeeper/configuration.xsl /alluxio-share/zookeeper/$myip/config
cp ${DIR}/zookeeper/log4j.properties /alluxio-share/zookeeper/$myip/config

docker rm -f alluxio-zk
docker run -d \
  --name alluxio-zk \
  -e ZOO_MY_ID=${myid} \
  -e ZOO_SERVERS=${zoo_servers} \
  -p 2888:2888 \
  -p 3888:3888 \
  -p 2181:2181 \
  -v /alluxio-share/zookeeper/$myip/config:/conf \
  -v /alluxio-share/zookeeper/$myip/data:/data \
  -v /alluxio-share/zookeeper/$myip/log:/datalog \
  zookeeper:3.4
