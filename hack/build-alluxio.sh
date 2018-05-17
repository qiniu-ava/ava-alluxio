#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..

build_tarball=true
build_kodo=true
build_image=true

if [ `echo $@ | grep '\--tarball=false' | wc -l | sed -e 's/^[[:space:]]*//'` = "1" ]; then
  build_tarball=false
fi

if [ `echo $@ | grep '\--kodo=false' | wc -l | sed -e 's/^[[:space:]]*//'` = "1" ]; then
  build_kodo=false
fi

if [ `echo $@ | grep '\--image=false' | wc -l | sed -e 's/^[[:space:]]*//'` = "1" ]; then
  build_image=false
fi

########################################################################### 
# build alluxio tarball and extract
###########################################################################
if [ $build_tarball != "false" ]; then
  echo "building alluxio tarball"
  mkdir -p .tmp/alluxio .tmp/kodo .tmp/temp/kodo && \
    cd .tmp/alluxio && \
    $DIR/../alluxio/dev/scripts/generate-tarballs single && \
    tar xf alluxio-1.7.2-SNAPSHOT.tar.gz && \
    cd .. && \
    cp alluxio/alluxio-1.7.2-SNAPSHOT/lib/alluxio-underfs-oss-1.7.2-SNAPSHOT.jar ./temp/ && \
    cd temp/kodo/ && \
    jar xf ../alluxio-underfs-oss-1.7.2-SNAPSHOT.jar && \
    cd $DIR/.. && \
    echo -e "\n\n\n"
else
  echo -e "skip building alluxio tarball\n\n\n"
fi

########################################################################### 
# build kodo sdk
###########################################################################
if [ $build_kodo != "false" ]; then
  echo "building alluxio kodo sdk"
  cd kodo
  mvn -DskipTests -Dlicense.skip=true compile install
  rm -rf $DIR/../.tmp/temp/kodo/com && cp -r target/classes/com ../.tmp/temp/kodo/com
  cd $DIR/../.tmp/temp/kodo && \
    rm -f alluxio-underfs-oss-1.7.2-SNAPSHOT.jar && \
    jar -cf alluxio-underfs-oss-1.7.2-SNAPSHOT.jar . && \
    mv ./alluxio-underfs-oss-1.7.2-SNAPSHOT.jar $DIR/../.tmp/alluxio/alluxio-1.7.2-SNAPSHOT/lib/ && \
    cd $DIR/.. && \
    cp dev/kodo-libs/* .tmp/alluxio/alluxio-1.7.2-SNAPSHOT/lib/ && \
    echo -e "\n\n\n"
else
  echo -e "skip building kodo sdk\n\n\n"
fi


########################################################################### 
# build docker image
###########################################################################
if [ $build_image != "false" ]; then
  echo "building docker image"
  docker build -t alluxio -f ./docker/Dockerfile.alluxio .tmp/alluxio/alluxio-1.7.2-SNAPSHOT && \
    echo -e "\n\n\n"
else
  echo -e "skip building docker image\n\n\n"
fi

