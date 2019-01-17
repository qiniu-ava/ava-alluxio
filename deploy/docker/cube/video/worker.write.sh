#!/bin/bash

######################################################################
# worker node list:
# jq42 jq66 jq67
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video
export NODE_LIST="jq42 jq66 jq67"

${DIR}/../template/worker.write.sh "$@"
