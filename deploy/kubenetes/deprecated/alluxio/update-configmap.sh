#!/bin/bash

role=$1

kubectl delete -n ava configmap alluxio-${role}-readonly-config
kubectl create -n ava configmap alluxio-${role}-readonly-config --from-file=ALLUXIO_CONFIG=./alluxio.properties.${role}
