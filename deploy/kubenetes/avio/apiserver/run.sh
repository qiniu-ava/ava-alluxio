#!/bin/bash

oldDir="$( pwd )"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$DIR" || return

if [[ ! -f "$( command -v jq )" ]]; then
  echo "dependency jq is not installed, plz install jq first"
  exit 1
fi

git_dirty=$(($(git status -s --porcelain | wc -l)))
if [ $git_dirty -ne 0 ]; then
  echo "plz commit your code and merge to develop branch of qiniu-ava/ava-alluxio before deploy apiserver"
  exit 1
fi

port=$( jq ".server.port" apiserver.conf )
tag=$(git rev-parse --short=7 HEAD)
sed "s/<APISERVER_PORT>/${port}/g" config.yml.template | sed "s/<REPOS_TAG>/${tag}/g" > config.yml

kubectl delete configmap avio-apiserver-config
kubectl create configmap avio-apiserver-config --from-file=apiserver.conf
kubectl apply -f config.yml

cd "$oldDir" || return
