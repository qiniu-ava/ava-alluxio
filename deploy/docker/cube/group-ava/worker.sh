#!/bin/bash

######################################################################
# worker node list:
# jq39
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=group-ava
export NODE_LIST="jq39 jq57"

${DIR}/../template/worker.sh "$@"
