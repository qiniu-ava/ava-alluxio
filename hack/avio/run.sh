#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

printUsage() {
  echo "usage: $DIR/run.sh [appName] where appName could be one of apiserver/executor/walker"
  exit 1
}

launch() {
  cd "$DIR/../../tools/golang/src/qiniu.com/app" || exit
  env GOOS=darwin go install -tags debug ./...
  cd "$DIR/../.." || exit

  if [ "$1" == all ];then
    killall avio-apiserver
    nohup "$DIR/../../tools/golang/bin/avio-apiserver" -f "$DIR/avio-apiserver.conf" 2>&1 &
    killall avio-walker
    nohup "$DIR/../../tools/golang/bin/avio-walker" -f "$DIR/avio-walker.conf" 2>&1 &
    killall avio-executor
    nohup "$DIR/../../tools/golang/bin/avio-executor" -f "$DIR/avio-executor.conf" 2>&1 &
  else
    killall avio-"$1"
    "$DIR/../../tools/golang/bin/avio-$1" -f "$DIR/avio-$1.conf"
  fi
}

if [ "$#" -gt 1 ]; then
  printUsage
fi

app=$1

if [ "$app" == "" ];then
  app=all
elif [[ "$app" != "apiserver" && "$app" != "walker" && "$app" != "executor" ]]; then
  printUsage
fi

launch $app
