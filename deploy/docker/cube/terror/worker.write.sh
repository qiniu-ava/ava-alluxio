#!/bin/bash

######################################################################
# worker node list:
# jq69 jq71
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=terror
export NODE_LIST="jq69 jq71"

${DIR}/../template/worker.write.sh "$@"
