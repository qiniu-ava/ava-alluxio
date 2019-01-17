#!/bin/bash

######################################################################
# worker node list:
# jq66
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=blademaster
export NODE_LIST="jq66"

${DIR}/../template/worker.write.sh "$@"
