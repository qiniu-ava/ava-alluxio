#!/bin/bash

print_usage() {
  echo "usage ========="
  echo "$<REPO_ROOT>/docker/generate-image.sh <app>"
  echo "where <app> should be one of all/dev/alluxio/alluxio-dev/client/client-cpu/client-gpu/zookeeper or empty which equal to all"
}

# usage: str_index_of <str> <substr>
str_index_of() {
  first_part=`echo $1 | awk -F"$2" '{print $1}'`
  next_part=`echo $1 | awk -F"$2" '{print $2}'`
  if [ "$next_part" != "" ]; then
    count=`echo $first_part | wc -c`
    echo `expr $count - 1`
  else
    if [ "$first_part" != $1 ]; then
      count=`echo $first_part | wc -c`
      echo `expr $count - 1`
    else
      echo -1
    fi
  fi
}

build_alluxio() {
  echo ""
}

build_kodo_sdk() {
  echo ""
}

build_alluxio_prod() {
  echo ""
}

build_alluxio_dev() {
  echo ""
}

docker_build() {
  dot_index=`str_index_of $1 .`
  docker_file_name="Dockerfile.$1"
  declare image_name
  declare image_tag
  if [ "$dot_index" != "-1" ]; then
    image_name=`echo $1 | awk -F"." '{print $1}'`
    image_tag=`echo $1 | awk -F"." '{print $2}'`
  else
    image_name=$1
    image_tag="latest"
  fi

  echo "docker build -t $image_name:$image_tag -f ./$docker_file_name ../"
  docker build -t $image_name:$image_tag -f ./$docker_file_name ../
}

app=$1

main() {
  build_alluxio_flag=true
  build_kodo_sdk_flag=true

  if [ "$build_alluxio_flag" != "true" ];then
    build_alluxio
  fi

  if [ "$build_kodo_sdk_flag" != "true" ];then
    build_kodo_sdk
  fi

  if [ "$app" = "" ]; then
    app=all
  fi

  case ${app} in
    all)
      docker_build alluxio
      docker_build alluxio.dev
      docker_build caffe.cpu
      docker_build caffe.gpu
      docker_build client.dev
      docker_build zookeeper
      ;;
    alluxio)
      docker_build alluxio
      docker_build alluxio.dev
      docker_build client.dev
      docker_build zookeeper
      ;;
    dev)
      docker_build alluxio.dev
      docker_build client.dev
      ;;
    alluxio-dev)
      docker_build alluxio.dev
      ;;
    zookeeper)
      docker_build zookeeper
      ;;
    client)
      docker_build caffe.cpu
      docker_build caffe.gpu
      ;;
    client-cpu)
      docker_build caffe.cpu
      ;;
    client-gpu)
      docker_build caffe.gpu
      ;;
    *)
      print_usage
      exit 1
      ;;
  esac
}

main $@
