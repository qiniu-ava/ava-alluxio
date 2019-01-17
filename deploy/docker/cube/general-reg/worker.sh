#!/bin/bash

######################################################################
# worker node list:
# jq65 jq66 jq67 jq68 jq69 jq70 jq71 jq73
######################################################################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export GROUP=general-reg
export NODE_LIST="jq65 jq66 jq67 jq68 jq69 jq70 jq71 jq73"

${DIR}/../template/worker.sh "$@"
