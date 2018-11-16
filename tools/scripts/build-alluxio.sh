#!/bin/bash

###########################################################################
################################# Usage ###################################
## options:                                                              ##
##   -t/--tarball=true/false, whether to build alluxio tarball           ##
##      optional, default: true                                          ##
##   -k/--kodo=true/false whether to build kodo jar package              ##
##      optional, default: true                                          ##
##   -i/--image=true/false whether to build docker image                 ##
##      optional, default: true                                          ##
##   --local-alluxio=<absolute_path_to_your_alluxio_repostory>           ##
##      optional, default: alluxio submodule path                        ##
##   --local-kodo=<absolute_path_to_your_kodo_repostory>                 ##
##      optional, default: kodo submodule path                           ##
##   -p/--push=true/false whether to push docker image                   ##
##      optional, default: true                                          ##
###########################################################################

set -e

ALLUXIO_VERSION="1.8.1-SNAPSHOT"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../..

build_tarball=true
build_kodo=true
build_image=true
local_alluxio=$DIR/../../alluxio
local_kodo=$DIR/../../kodo
push_image=true

for i in "$@"; do
  case $i in
    -t=*|--tarball=*)
      build_tarball="${i#*=}"
    ;;
    -k=*|--kodo=*)
      build_kodo="${i#*=}"
    ;;
    -i=*|--image=*)
      build_image="${i#*=}"
    ;;
    --local-alluxio=*)
      local_alluxio="${i#*=}"
    ;;
    --local-kodo=*)
      local_kodo="${i#*=}"
    ;;
    -p=*|--push=*)
      push_image="${i#*=}"
    ;;
    *)
      # unknown option
    ;;
  esac
done

mkdir -p .tmp/alluxio .tmp/kodo .tmp/temp/kodo

alluxio_version=$(grep -m1 "<version>" ${local_alluxio}/pom.xml | awk -F'>|<' '{print $3}')

########################################################################### 
# build alluxio tarball and extract
###########################################################################
cd $DIR/../..
if [ $build_tarball != "false" ]; then
  echo "building alluxio tarball"
  rm -rf .tmp/alluxio/* .tmp/temp/kodo/* && \
    cd .tmp/alluxio && \
    ${local_alluxio}/dev/scripts/generate-tarballs single && \
    tar xf alluxio-"${alluxio_version}".tar.gz && \
    rm alluxio-"${alluxio_version}".tar.gz && \
    cd .. && \
    cp alluxio/alluxio-"${alluxio_version}"/lib/alluxio-underfs-oss-${ALLUXIO_VERSION}.jar ./temp/ && \
    cd temp/kodo/ && \
    jar xf ../alluxio-underfs-oss-"${alluxio_version}".jar
  cd $DIR/../.. && \
    echo -e "\n\n\n"
else
  echo -e "skip building alluxio tarball\n\n\n"
fi

########################################################################### 
# build kodo sdk
###########################################################################
cd $DIR/../..
if [ $build_kodo != "false" ]; then
  echo "building alluxio kodo sdk"
  cd ${local_kodo}
  mvn -DskipTests -Dlicense.skip=true compile install
  rm -rf $DIR/../../.tmp/temp/kodo/com && cp -r target/classes/com $DIR/../../.tmp/temp/kodo/com
  cd $DIR/../../.tmp/temp/kodo && \
    rm -f alluxio-underfs-oss-"${alluxio_version}".jar && \
    jar -cf alluxio-underfs-oss-"${alluxio_version}".jar . && \
    mv ./alluxio-underfs-oss-"${alluxio_version}".jar $DIR/../../.tmp/alluxio/alluxio-${ALLUXIO_VERSION}/lib/ && \
    cd $DIR/../.. && \
    cp $DIR/docker-image/kodo-libs/* .tmp/alluxio/alluxio-"${alluxio_version}"/lib/ && \
    echo -e "\n\n\n"
else
  echo -e "skip building kodo sdk\n\n\n"
fi

###########################################################################
# build alluxio client tarball
###########################################################################
cd $DIR/../..
cd .tmp/alluxio/alluxio-${ALLUXIO_VERSION}
cp ../../../deploy/env/alluxio-flex-volume.sh ./ && cp ../../../deploy/env/client/alluxio-* ./conf
cd ../
if git describe --exact-match --tags $(git rev-parse --short=7 HEAD); then
  tag=`git describe --exact-match --tags $(git rev-parse --short=7 HEAD)` && tar zcvf ${tag}.tar.gz ./alluxio-${ALLUXIO_VERSION}
else
  tag=`git rev-parse --short=7 HEAD` && tar zcvf ava-alluxio-${tag}.tar.gz ./alluxio-${ALLUXIO_VERSION}
fi
cd $DIR/../..

########################################################################### 
# build docker image
###########################################################################
cd $DIR/../..
if [ $build_image != "false" ]; then
  echo "building docker image"
  cp $DIR/docker-image/Dockerfile.alluxio $DIR/docker-image/entrypoint.sh .tmp/alluxio/
  cd $local_alluxio && alluxio_hash=`git rev-parse --short=7 HEAD` && cd -
  cd $local_kodo && kodo_hash=`git rev-parse --short=7 HEAD` && cd -
  docker build -t alluxio:$alluxio_hash-$kodo_hash --build-arg ALLUXIO_VERSION="${alluxio_version}" -f .tmp/alluxio/Dockerfile.alluxio .tmp/alluxio
  if [ $push_image != "false" ]; then
    docker tag alluxio:$alluxio_hash-$kodo_hash reg-xs.qiniu.io/atlab/alluxio:$alluxio_hash-$kodo_hash
    docker push reg-xs.qiniu.io/atlab/alluxio:$alluxio_hash-$kodo_hash
  fi
  echo -e "\n\n\n"
else
  echo -e "skip building docker image\n\n\n"
fi
