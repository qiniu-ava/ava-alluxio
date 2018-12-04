#!/bin/bash

######################################################################
# master node list:
# jq67
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=terror
export NODE_LIST="jq67"

${DIR}/../template/master.sh "$@"
