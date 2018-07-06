#!/bin/bash

export LOGKIT_VERSION=v1.5.0
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR
if [ ! -d "./_package/" ];then
  wget https://pandora-dl.qiniu.com/logkit_${LOGKIT_VERSION}.tar.gz && \
  tar xvf logkit_${LOGKIT_VERSION}.tar.gz && \
  rm logkit_${LOGKIT_VERSION}.tar.gz && cd _package/
  rm logkit.conf
  mkdir confs && mv ../confs/* confs && rm -r ../confs
  mkdir script && mv ../script/* script && rm -r ../script
  mv ../logkit.sh ./
  mv ../logkit.conf ./
else
  echo "logkit has been installed"
fi

scriptDIR=$DIR/_package/script
for file in $DIR/_package/confs/* ; do
  temp_file=`basename $file`
  sh_name=`basename $file | cut -d '.' -f 2 | sed 's/_/-/'`
  scriptName=\"$scriptDIR/readLog-$sh_name.sh\"
  sed -i 's!\"log_path\": <log_path>!\"log_path\": '$scriptName'!' $DIR/_package/confs/$temp_file
done
