#!/bin/bash

export LOGKIT_VERSION=v1.5.1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR" || exit
if [ ! -d "$HOME/logkit/_package/" ];then
  if [ ! -d "$HOME/logkit/" ];then
    mkdir "$HOME"/logkit
  fi
  wget https://pandora-dl.qiniu.com/logkit_${LOGKIT_VERSION}.tar.gz && \
  tar xvf logkit_${LOGKIT_VERSION}.tar.gz && \
  rm logkit_${LOGKIT_VERSION}.tar.gz && \
  mv _package ~/logkit/
  rm ~/logkit/_package/logkit.conf
  mkdir ~/logkit/_package/confs && cp ./confs/* ~/logkit/_package/confs
  mkdir ~/logkit/_package/script && cp ./script/* ~/logkit/_package/script
  cp ./start.sh ~/logkit/_package/
  cp ./logkit.conf ~/logkit/_package/
  cd ~/logkit/_package/ || exit
  scriptDIR=$PWD/script
  for file in "$PWD"/confs/* ; do
    temp_file=$(basename "$file")
    sh_name=$(basename "$file" | cut -d '.' -f 2 | sed 's/_/-/')
    scriptName=\"$scriptDIR/readLog-$sh_name.sh\"
    sed -i 's!\"log_path\": <log_path>!\"log_path\": '"$scriptName"'!' "$PWD"/confs/"$temp_file"
  done
else
  echo "logkit has been installed. If you want to reinstall logkit please rm -r $HOME/logkit/_package/"
fi
