#!/bin/bash

export LOGKIT_VERSION=v1.2.0
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR" || exit
if [ ! -d "$HOME/logkit/_package/" ];then
  if [ ! -d "$HOME/logkit/" ];then
    mkdir "$HOME"/logkit
  fi
  wget https://pandora-dl.qiniu.com/logkit-pro-local_linux64_${LOGKIT_VERSION}.tar.gz && \
  tar xzvf logkit-pro-local_linux64_${LOGKIT_VERSION}.tar.gz && \
  rm logkit-pro-local_linux64_${LOGKIT_VERSION}.tar.gz && \
  mv logkit-pro_linux64_${LOGKIT_VERSION} ~/logkit/
  mkdir ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/confs && cp ./confs/* ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/confs
  mkdir ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/script && cp ./script/* ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/script
  cp ./start.sh ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/
  cp ./logkit.conf ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/
  cd ~/logkit/logkit-pro_linux64_${LOGKIT_VERSION}/ || exit
  scriptDIR=$PWD/script
  for file in "$PWD"/confs/* ; do
    temp_file=$(basename "$file")
    sh_name=$(basename "$file" | cut -d '.' -f 2 | sed 's/_/-/')
    scriptName=\"$scriptDIR/readLog-$sh_name.sh\"
    sed -i 's!\"log_path\": \"<script_path>\"!\"log_path\": '"$scriptName"'!' "$PWD"/confs/"$temp_file"
  done
else
  echo "logkit has been installed. If you want to reinstall logkit please rm -r $HOME/logkit/_package/"
fi
