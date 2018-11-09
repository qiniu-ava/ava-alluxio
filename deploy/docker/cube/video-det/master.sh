#!/bin/bash

######################################################################
# master node list:
# jq54
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video-det
export NODE_LIST="jq54"

${DIR}/../template/master.sh "$@"
