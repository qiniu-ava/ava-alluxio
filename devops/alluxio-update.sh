#!/bin/bash

set -ex

alluxio_tag=$1

if [ -z $alluxio_tag ]; then
  echo "alluxio tag not set"
  exit 1
fi

cd /opt/
curl -s "http://atlab.ava-test.xbowen.com/release/alluxio/ava-alluxio-$alluxio_tag.tar.gz" -O

tar zxf "./ava-alluxio-$alluxio_tag.tar.gz"
if [ -d ./alluxio ]; then
  mv ./alluxio ./alluxio-$(date -Iseconds)
fi

mv ./alluxio-1.8.1-SNAPSHOT ./alluxio

mkdir -p /var/log/alluxio/

