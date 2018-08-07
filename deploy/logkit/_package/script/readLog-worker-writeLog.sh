#!/bin/bash
logTime=2018-07-01T00:00:00.689381000Z #SAVE_TIME
dockerName="alluxio-worker-writer"
dockerID=$(docker ps --no-trunc | awk '$NF=="alluxio-worker-writer"{print $1}')

if [ "$dockerID" != "" ]; then
    newTime=$(docker logs --timestamps --tail 1 "$dockerID" | awk '{print $1}')
    if [ "$logTime" != "$newTime" ]; then
        logs=$(docker logs --since="$logTime" "$dockerID" | awk 'NR%1==0{print;next}')
        if [ "$logs" != "" ]; then
            echo $dockerName"|""$dockerID""|""$logs"
            sed -i "1,3s/.*#SAVE_TIME/logTime=$newTime #SAVE_TIME/" "$0"
        fi
    fi
fi
