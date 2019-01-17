#!/bin/bash

######################################################################
# worker node list:
# jq65
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=terror
export NODE_LIST="jq65"

${DIR}/../template/worker.write.sh "$@"
