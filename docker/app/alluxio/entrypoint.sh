#!/bin/bash

set -ex

service nginx start

if [ $? -eq 0 ]; then
  sleep infinity
else
  echo "start nginx service failed"
  echo -e "\n\n\nerror logs:\n"
  cat /var/log/nginx/error.log
  exit 1
fi

