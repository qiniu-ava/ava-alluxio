#!/bin/bash

######################################################################
# worker node list:
# jq54 jq56 jq57
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video-det
export NODE_LIST="jq54 jq56 jq57"

${DIR}/../template/worker.sh "$@"
