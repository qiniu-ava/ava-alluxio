#!/bin/bash

###########################################################################
################################# Usage ###################################
## options:                                                              ##
##   -t/--tarball=true/false, whether to build alluxio tarball           ##
##      optional, default: true                                          ##
##   -i/--image=true/false whether to build docker image                 ##
##      optional, default: true                                          ##
##   --local-alluxio=<absolute_path_to_your_alluxio_repostory>           ##
##      optional, default: alluxio submodule path                        ##
##   -p/--push=true/false whether to push docker image                   ##
##      optional, default: true                                          ##
###########################################################################

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../..

build_tarball=true
build_image=true
local_alluxio=$DIR/../../alluxio
push_image=true

for i in "$@"; do
  case $i in
    -t=*|--tarball=*)
      build_tarball="${i#*=}"
    ;;
    -i=*|--image=*)
      build_image="${i#*=}"
    ;;
    --local-alluxio=*)
      local_alluxio="${i#*=}"
    ;;
    -p=*|--push=*)
      push_image="${i#*=}"
    ;;
    *)
      # unknown option
    ;;
  esac
done


ALLUXIO_VERSION=$(grep -m1 "<version>" ${local_alluxio}/pom.xml | awk -F'>|<' '{print $3}')

mkdir -p .tmp/alluxio

########################################################################### 
# build alluxio tarball
###########################################################################
cd $DIR/../..
if [ $build_tarball != "false" ]; then
  echo "building alluxio tarball"
  rm -rf .tmp/alluxio/* && \
    cd .tmp/alluxio && \
    ${local_alluxio}/dev/scripts/generate-tarballs single && \
    tar xf alluxio-"${ALLUXIO_VERSION}".tar.gz && \
    rm alluxio-"${ALLUXIO_VERSION}".tar.gz && \
    cd $DIR/../../.tmp/alluxio/alluxio-${ALLUXIO_VERSION} && \
    cp $DIR/../../deploy/env/alluxio-flex-volume.sh ./ && \
    cp $DIR/../../deploy/env/client/alluxio-* ./conf
  cd ..
  if git describe --exact-match --tags $(git rev-parse --short=7 HEAD); then
    tag=`git describe --exact-match --tags $(git rev-parse --short=7 HEAD)` && \
      echo "$tag" > ./alluxio-${ALLUXIO_VERSION}/version && \
      tar zcvf ${tag}.tar.gz ./alluxio-${ALLUXIO_VERSION}
  else
    tag=`git rev-parse --short=7 HEAD` && \
      echo "ava-alluxio-$tag" > ./alluxio-${ALLUXIO_VERSION}/version && \
      tar zcvf ava-alluxio-${tag}.tar.gz ./alluxio-${ALLUXIO_VERSION}
  fi
  cd $DIR/../.. && \
    echo -e "\n\n\n"
else
  echo -e "skip building alluxio tarball\n\n\n"
fi

########################################################################### 
# build docker image
###########################################################################
cd $DIR/../..
if [ $build_image != "false" ]; then
  echo "building docker image"
  cp $DIR/docker-image/Dockerfile.alluxio $DIR/docker-image/entrypoint.sh .tmp/alluxio/
  cd $local_alluxio && alluxio_hash=`git rev-parse --short=7 HEAD` && cd -
  kodo_hash="kodo"
  docker build -t alluxio:$alluxio_hash-$kodo_hash --build-arg ALLUXIO_VERSION="${ALLUXIO_VERSION}" -f .tmp/alluxio/Dockerfile.alluxio .tmp/alluxio
  if [ $push_image != "false" ]; then
    docker tag alluxio:$alluxio_hash-$kodo_hash reg-xs.qiniu.io/atlab/alluxio:$alluxio_hash-$kodo_hash
    docker push reg-xs.qiniu.io/atlab/alluxio:$alluxio_hash-$kodo_hash
  fi
  echo -e "\n\n\n"
else
  echo -e "skip building docker image\n\n\n"
fi
