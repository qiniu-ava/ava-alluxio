#!/bin/bash
logTime=2018-07-01T06:47:36.103266200Z #SAVE_TIME

dockerID=$(docker ps --no-trunc | awk '$NF=="<dockerName>"{print $1}')

if [ "$dockerID" != "" ]; then
    newTime=$(docker logs --timestamps --tail 1 "$dockerID" | awk '{print $1}')
    if [ "$logTime" != "$newTime" ]; then
        logs=$(docker logs --since="$logTime" "$dockerID")
        if [ "$logs" != "" ]; then
            if [[ $(wc -l "$logs") -gt 200 ]]; then
                log_front=$(echo "$logs" | head -100)
                log_end=$(echo "$logs" | tail -100)
                logs="$log_front""\\n""$log_end"
            fi
            echo "$dockerID""|""$logs""|""<dockerName>"
            sed -i "1,3s/.*#SAVE_TIME/logTime=$newTime #SAVE_TIME/" "$0"
        fi
    fi
fi
