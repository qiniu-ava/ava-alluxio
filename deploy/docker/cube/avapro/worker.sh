#!/bin/bash

######################################################################
# worker node list:
# alluxio-test-fusion
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ALLUXIO_ENV=/data/alluxio/env

export GROUP=avapro
export NODE_LIST="alluxio-test-fusion"

${DIR}/../template/worker.sh "$@"
