#!/bin/bash

######################################################################
# master node list:
# jw3 jw10 jw11
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=juewa-default
export NODE_LIST="jw3 jw10 jw11"

${DIR}/../template/master.sh "$@"
