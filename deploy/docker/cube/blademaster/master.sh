#!/bin/bash

######################################################################
# master node list:
# jq71
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=blademaster
export NODE_LIST="jq71"

${DIR}/../template/master.sh "$@"
