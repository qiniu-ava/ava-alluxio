#!/bin/bash

######################################################################
# master node list:
# jq70
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=terror
export NODE_LIST="jq70"

${DIR}/../template/master.sh "$@"
