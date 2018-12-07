#!/bin/bash

######################################################################
# worker node list:
# jq68 jq69
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=ocr
export NODE_LIST="jq68 jq69"

${DIR}/../template/worker.sh "$@"
