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
###########################################################################

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..

build_tarball=true
build_kodo=true
build_image=true
local_alluxio=$DIR/../alluxio
local_kodo=$DIR/../kodo

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
    *)
      # unknown option
    ;;
  esac
done

mkdir -p .tmp/alluxio .tmp/kodo .tmp/temp/kodo

########################################################################### 
# build alluxio tarball and extract
###########################################################################
if [ $build_tarball != "false" ]; then
  echo "building alluxio tarball"
  rm -rf .tmp/alluxio/* .tmp/temp/kodo/* && \
    cd .tmp/alluxio && \
    ${local_alluxio}/dev/scripts/generate-tarballs single && \
    tar xf alluxio-1.7.2-SNAPSHOT.tar.gz && \
    rm alluxio-1.7.2-SNAPSHOT.tar.gz && \
    cd .. && \
    cp alluxio/alluxio-1.7.2-SNAPSHOT/lib/alluxio-underfs-oss-1.7.2-SNAPSHOT.jar ./temp/ && \
    cd temp/kodo/ && \
    jar xf ../alluxio-underfs-oss-1.7.2-SNAPSHOT.jar
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
  cd ${local_kodo}
  mvn -DskipTests -Dlicense.skip=true compile install
  rm -rf $DIR/../.tmp/temp/kodo/com && cp -r target/classes/com $DIR/../.tmp/temp/kodo/com
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
  cp ./docker/entrypoint.sh .tmp/alluxio/
  ls .tmp/alluxio
  docker build -t alluxio -f ./docker/Dockerfile.alluxio .tmp/alluxio
  echo -e "\n\n\n"
else
  echo -e "skip building docker image\n\n\n"
fi
