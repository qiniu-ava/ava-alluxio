#!/bin/bash

############################################################
################################# Usage ####################
## options:                                               ##
##   --ak=string, pandora account AK                      ##
##      optional, default: currentPath account AK         ##
##   --sk=string, pandora account SK                      ##
##      optional, default: currentPath account SK         ##
##   --aksk=string, pandora account absolute path         ##
##      optional, default: currentPath/account            ##
############################################################
##account file should contain:                            ##
##                                                        ##
##  AK=XXXXXXXXXXX                                        ##
##  SK=XXXXXXXXXXX                                        ##
##                                                        ##
############################################################

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

AK="$(sed -n 's/AK=//gp' $DIR/account)"
SK="$(sed -n 's/SK=//gp' $DIR/account)"

for i in "$@"; do
  case $i in
    --ak=*)
      AK="${i#*=}"
    ;;
    --sk=*)
      SK="${i#*=}"
    ;;
    --AKSK=*)
	  akskFile="${i#*=}"
      AK="$(sed -n 's/AK=//gp' $akskFile)"
      SK="$(sed -n 's/SK=//gp' $akskFile)"
    ;;
    *)
      # unknown option
    ;;
  esac
done

if [ "$AK" != "" -a "$SK" != "" ]; then
  for file in $DIR/confs/* ; do  
    temp_file=`basename $file`
    sed -i "" 's/\"pandora_ak\": <pandora_ak>/\"pandora_ak\": '$AK'/' $DIR/confs/$temp_file
    sed -i "" 's/\"pandora_sk\": <pandora_sk>/\"pandora_ak\": '$SK'/' $DIR/confs/$temp_file
  done
  ./logkit -f logkit.conf
else
  echo "please check your AKSK or account path"
fi
  
