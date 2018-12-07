#!/bin/bash

######################################################################
# master node list:
# jq73
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=ocr
export NODE_LIST="jq73"

${DIR}/../template/master.sh "$@"
