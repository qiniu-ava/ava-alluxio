#!/bin/bash

reborn_master() {
  kubectl delete -f ./master.yml
  kubectl delete statefulset/alluxio-master-bowen --cascade=false
  kubectl delete pod/alluxio-master-bowen-0
  kubectl delete pod/alluxio-master-bowen-1
  kubectl delete pod/alluxio-master-bowen-2
}

reborn_worker() {
  kubectl delete -f ./worker.yml
}

case ${1} in
  all)
    reborn_master
    reborn_worker
    ;;
  master)
    reborn_master
    ;;
  worker)
    reborn_worker
    ;;
  *)
    echo "usage ./reborn.sh <master/worker/all>"
    ;;
esac

