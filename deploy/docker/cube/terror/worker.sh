#!/bin/bash

######################################################################
# worker node list:
# jq65 jq66 jq67
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=terror
export NODE_LIST="jq65 jq66 jq67"

${DIR}/../template/worker.sh "$@"
