#!/bin/bash

if [ $# -lt 2 ]
then
  echo usage "$0 <alluxiosc.so path> <list file>"
  exit -1
fi

so=$1
list=$2
if [ ! -e $so ]
then
  echo "so file $so does not exist"
  exit -1
fi
if [ ! -e $list ]
then
  echo "list file $list does not exist"
  exit -1
fi

MAX=50

cache=()
trim() {
  left=()
  for f in ${cache[@]}
  do
    if [ ! -e "$f" ]
    then
      echo "remove non existing $f ..."
      continue
    fi
    ok=`alluxiosc_query=1 LD_PRELOAD=$so ls "$f" | grep 'query=/'`
    if [ -z "$ok" ]
    then
      echo "trying to cache again $f ..."
      left+=("$f")
      alluxiosc_cache=1 LD_PRELOAD=$so ls "$f" > /dev/null &
    else
      echo "done cache $f ..."
    fi
  done
  cache=(${left[@]})
}

while read line 
do
  while [ ${#cache[@]} -gt $MAX ]
  do
    sleep 2
    trim
  done

  f=`echo $line | awk '{print $1}'`
  cache+=("$f")
  echo "trying to cache $f ..."
  alluxiosc_cache=1 LD_PRELOAD=$so ls "$f" > /dev/null &
done < $list

while [ ${#cache[@]} -gt 0 ]
do
  sleep 1
  trim
done

