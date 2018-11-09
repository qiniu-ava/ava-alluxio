#!/bin/bash

######################################################################
# worker node list:
# jq40 jq41 jq42 jq56 jq57
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video
export NODE_LIST="jq40 jq41 jq42 jq56 jq57"

${DIR}/../template/worker.sh "$@"
