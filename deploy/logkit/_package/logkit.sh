#!/bin/bash

############################################################
################################# Usage ####################
## options:                                               ##
##   --ak=string, pandora account AK                      ##
##      optional, default: null                           ##
##   --sk=string, pandora account SK                      ##
##      optional, default: null                           ##
##   --aksk=string, pandora aksk absolute path            ##
##      optional, default: null                           ##
##   --mail=string, default:ava-dev                       ##
##      optional, mail with pandora service               ##
############################################################
##account file should contain:                            ##
##                                                        ##
##  mail_AK=XXXXXXXXXXX                                   ##
##  mail_SK=XXXXXXXXXXX                                   ##
##                                                        ##
############################################################

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

AK=""
SK=""
service_mail="ava-dev"

for i in "$@"; do
  case $i in
    --ak=*)
      AK="${i#*=}"
    ;;
    --sk=*)
      SK="${i#*=}"
    ;;
    --mail=*)
      service_mail="${i#*=}"
    ;;
    --aksk=*)
	  akskFile="${i#*=}"
      AK="$(sed -n 's/'$service_mail'_AK=//gp' $akskFile)"
      SK="$(sed -n 's/'$service_mail'_SK=//gp' $akskFile)"
    ;;
    *)
      # unknown option
    ;;
  esac
done

if [ "$AK" != "" -a "$SK" != "" ]; then
  for file in $DIR/confs/* ; do
    temp_file=`basename $file`
    sed -i 's/\"pandora_ak\": <pandora_ak>/\"pandora_ak\": '$AK'/' $DIR/confs/$temp_file
    sed -i 's/\"pandora_sk\": <pandora_sk>/\"pandora_sk\": '$SK'/' $DIR/confs/$temp_file
  done
  logkitPID="$(ps -aux | grep "logkit" | grep -v grep | awk '{print $2}')"
  if [ "$logkitPID" != "" ]; then
    for i in $logkitPID; do
      kill -9 $logkitPID
    done
  fi
  nohup ./logkit -f logkit.conf > ../logkit.out 2>&1 &
else
  echo "please check your AKSK or account path"
fi
