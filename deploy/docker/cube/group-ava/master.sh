#!/bin/bash

######################################################################
# master node list:
# jq15
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=group-ava
export NODE_LIST="jq15"

${DIR}/../template/master.sh "$@"
