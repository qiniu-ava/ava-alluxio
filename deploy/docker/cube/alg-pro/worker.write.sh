#!/bin/bash

######################################################################
# worker node list:
# jq20
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=alg-pro
export NODE_LIST="jq20"

${DIR}/../template/worker.write.sh "$@"
