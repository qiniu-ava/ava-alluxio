#!/bin/bash

oldDir=$(pwd)
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$DIR" || return

push=false
for i in "$@"; do
  case $i in
    -p|--push)
      push=true
    ;;
  esac
done

if [ "$push" == true ];then
  git_dirty=$(($(git status -s --porcelain | wc -l)))
  if [[ $git_dirty -ne 0 ]];then
    echo "plz commit your code and merge to develop branch of qiniu-ava/ava-alluxio before build docker images"
    exit 1
  fi
fi

tag=$(git describe --exact-match --tags "$(git rev-parse --short=7 HEAD)" 2>&1)
if echo "$tag" | grep -qE "^fatal: no tag"; then
  tag=$(git rev-parse --short=7 HEAD)
fi

build_avio_app() {
  app="$1"
  if [ "$app" == all ]; then
    build_avio_app apiserver
    build_avio_app executor
    build_avio_app walker
  elif [ "$app" == "apiserver" ] || [ "$app" == "executor" ] || [ "$app" == "walker" ];then
    cd "$DIR" || return
    docker build -t "reg-xs.qiniu.io/atlab/avio-$app:$tag" -f "$DIR/docker/app/avio-$app/Dockerfile" .
    if [ "$push" == true ];then
      docker push "reg-xs.qiniu.io/atlab/avio-$app:$tag"
    fi
    cd "$oldDir" || return
  else
    echo "invalid avio app name $app"
  fi
}

build_kafka() {
  cd "$DIR"/docker/app/kafka || return
  docker build -t reg-xs.qiniu.io/atlab/avio-kafka:"$tag" .
  if [ "$push" == true ]; then
    docker push reg-xs.qiniu.io/atlab/avio-kafka:"$tag"
  fi
  cd "$DIR" || return
}

build_zookeeper() {
  cd "$DIR"/docker/app/zookeeper || return
  docker build -t reg-xs.qiniu.io/atlab/avio-zookeeper:"$tag" .
  if [ "$push" == true ]; then
    docker push reg-xs.qiniu.io/atlab/avio-zookeeper:"$tag"
  fi
  cd "$DIR" || return
}

build_alluxio_dashboard() {
  cd "$DIR"/docker/app/alluxio || return
  docker build -t reg-xs.qiniu.io/atlab/alluxio-nginx:"$tag" .
  if [ "$push" == true ]; then
    docker push reg-xs.qiniu.io/atlab/alluxio-nginx:"$tag"
  fi
  cd "$DIR" || return
}

if [ "${#}" -eq 0 ] || [ "$1" == all ]; then
  # build all apps
  build_avio_app all
  build_kafka
  build_zookeeper
  build_alluxio_dashboard
else
  for i in "$@"; do
    case $i in
      avio)
        build_avio_app all
      ;;
      kafka)
        build_kafka
      ;;
      zookeeper)
        build_zookeeper
      ;;
      dashboard)
        build_alluxio_dashboard
      ;;
      -p|--push)
      ;;
      *)
        echo "usage: ${BASH_SOURCE[0]} [-p|--push] [appName] where appName could be one of all/avio/kafka/zookeeper/dashboard"
      ;;
    esac
  done
fi

cd "$oldDir" || return
