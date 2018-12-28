#!/bin/bash

######################################################################
# worker node list:
# jq40 jq41 jq42
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=face
export NODE_LIST="jq40 jq41 jq42"

${DIR}/../template/worker.write.sh "$@"
