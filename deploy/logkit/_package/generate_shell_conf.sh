#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )" || exit
if [ ! -d "$PWD/script" ];then
  mkdir script
fi
if [ ! -d "$PWD/confs" ];then
  mkdir confs
fi

generate_script() {
  cp "$PWD"/template/template.sh "$PWD"/script 
  mv "$PWD"/script/template.sh "$PWD"/script/readLog-"$1"Log.sh
  sed -i 's/\"<dockerName>\"/\"'"$1"'\"/' "$PWD"/script/readLog-"$1"Log.sh
}

genreate_conf() {
  cp "$PWD"/template/template.conf "$PWD"/confs
  mv "$PWD"/confs/template.conf "$PWD"/confs/runner."$1"Log.conf
  sed -i 's/<dockerName>/'"$1"'/' "$PWD"/confs/runner."$1"Log.conf
  sed -i 's/\"<workflow_name>\"/\"'"$1"'\"/' "$PWD"/confs/runner."$1"Log.conf
  sed -i 's/\"<repo_name>\"/\"'"$1"'\"/' "$PWD"/confs/runner."$1"Log.conf
}



for docker in $(docker ps | awk '/alluxio-worker/||/alluxio-master/{print $NF}'); do
  generate_script "$docker"
  genreate_conf "$docker"
done
