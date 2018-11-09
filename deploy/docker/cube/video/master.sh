#!/bin/bash

######################################################################
# master node list:
# jq39
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video
export NODE_LIST="jq39"

${DIR}/../template/master.sh "$@"
