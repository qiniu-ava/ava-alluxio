#!/bin/bash

myip=`ifconfig | grep 'inet addr:10.' | awk -F':' '{print $2}' | awk '{print $1}'`

docker rm -f alluxio-zk
docker run -d \
  --name alluxio-zk \
  -e ZOO_MY_ID=${myip} \
  -e ZOO_SERVERS=10.200.20.91:2888:3888,10.200.20.70:2888:3888,10.200.20.80:2888:3888 \
  -p 2888:2888 \
  -p 3888:3888 \
  -p 2181:2181 \
  zookeeper:3.4
