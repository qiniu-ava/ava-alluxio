#!/bin/bash

######################################################################
# master node list:
# jq72
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=general-reg
export NODE_LIST="jq72"

${DIR}/../template/master.sh "$@"
