#!/bin/bash

set -e

echo "hello world"

my_ip=`ifconfig | grep 'inet addr:172' | awk -F':' '{print $2}' | awk '{print $1}'`

echo "my_ip: $my_ip"
my_id=`echo $my_ip | awk -F'.' '{print $4}'`
echo "my_id: $my_id"
export ZOO_MY_ID=$my_id

# /docker-entrypoint.sh

sleep 3600
