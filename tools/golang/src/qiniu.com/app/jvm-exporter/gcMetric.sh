#!/bin/sh

getMetric()
{
  pid="$(docker exec "$1" jps | awk '$2=="'"$2"'"{print $1}')"
  if [ "$pid" != "" ]; then 
    metrics="$(docker exec "$1" jstat -gcutil "${pid}" | awk 'END{for(i=7;i<=NF;i++)printf $i" "}')"
    echo "{\"$1\":\"$metrics\"}"
    # echo "java_young_gc_times{name=\"$1\"} `echo ${metrics} | cut -d" " -f1`"
    # echo "java_young_gc_duration{name=\"$1\"} `echo ${metrics} | cut -d" " -f2`"
    # echo "java_full_gc_times{name=\"$1\"} `echo ${metrics} | cut -d" " -f3`"
    # echo "java_full_gc_duration{name=\"$1\"} `echo ${metrics} | cut -d" " -f4`"
    # echo "java_all_gc_duration{name=\"$1\"} `echo ${metrics} | cut -d" " -f5`"
    

    for gcTID in $(docker exec "$1" jstack "${pid}" | grep -w '\<Parall.*GC\>' | sed 's/.*nid=\(.*\)runnable.*/\1/g'); do
      id=$(printf %d "$gcTID")
      hz=$(getconf CLK_TCK)
      cpuTime=$(docker exec "$1" cat /proc/"$pid"/task/"$id"/stat)
      utime=$(echo "$cpuTime" | cut -d" " -f14)
      stime=$(echo "$cpuTime" | cut -d" " -f15)
      cpuUse=$(awk 'BEGIN{print('$utime'/'$hz'+'$stime'/'$hz')}')
      echo "{\"$1 $id\":\"$cpuUse\"}"
    #   echo "java_gc_thread_cpu_used{name=\"$1\", threads=\"$id\"} $cpuUse"
    done
  fi
}

for docker in $(docker ps | awk '/alluxio/{print $NF}'); do
  case $docker in
    alluxio-master*)
      getMetric "$docker" "AlluxioMaster"
    ;;
    alluxio-worker*)
      getMetric "$docker" "AlluxioWorker"
    ;;
    *)
      # unknown option
    ;;
  esac
done
