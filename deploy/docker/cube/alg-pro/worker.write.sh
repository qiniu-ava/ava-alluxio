#!/bin/bash

######################################################################
# worker node list:
# jq19 jq20 jq21
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=alg-pro
export NODE_LIST="jq19 jq20 jq21"

${DIR}/../template/worker.write.sh "$@"
