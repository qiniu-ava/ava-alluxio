#!/bin/bash

######################################################################
# master node list:
# jq16
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=face
export NODE_LIST="jq16"

${DIR}/../template/master.sh "$@"
