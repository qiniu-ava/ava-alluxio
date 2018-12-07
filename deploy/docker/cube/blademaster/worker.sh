#!/bin/bash

######################################################################
# worker node list:
# jq65 jq66 jq67 jq68 jq69
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=blademaster
export NODE_LIST="jq65 jq66 jq67 jq68 jq69"

${DIR}/../template/worker.sh "$@"
