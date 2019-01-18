#!/bin/bash

######################################################################
# worker node list:
# jq42 jq65 jq69
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=video
export NODE_LIST="jq42 jq65 jq69"

${DIR}/../template/worker.write.sh "$@"
