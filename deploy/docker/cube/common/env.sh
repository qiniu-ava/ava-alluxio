#!/bin/bash

####################################################################
######## alluxio cluster environment settings include:      ########
######## 1. zookeeper service address                       ########
######## 2. cluster ethernet network gateway                ########
######## for example:                                       ########
######## cat ${ALLUXIO_ENV}/cluster                         ########
######## export ALLUXIO_ZOOKEEPER_ADDRESS=<zk_server_addr>  ########
######## export ALLUXIO_IP_PREFIX=192.168                   ########
######## export ALLUXIO_IP_EXCLUDE=192.168.212              ########
####################################################################

if [ "${ALLUXIO_ENV}" = "" ]; then
  ALLUXIO_ENV=/disk-cephfs/alluxio/env
fi

# default cluster setting, for jq
ALLUXIO_CLUSTER_NAME=jq-alluxio
ALLUXIO_ZOOKEEPER_ADDRESS=192.168.213.42:2181,192.168.213.45:2181,192.168.213.46:2181
ALLUXIO_IP_PREFIX=192.168.
ALLUXIO_IP_EXCLUDE=192.168.212
