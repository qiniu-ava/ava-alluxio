#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ ! -f "$DIR/env" ]]; then
  echo "plz config your setting in path $DIR/env, env.template may be helpful for you"
  exit 1
fi

# shellcheck disable=SC1090
source "$DIR"/env
myip=$(ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | awk '{print $2}')
sed "s/<HOST_NAME>/${myip}/g" "$DIR"/docker-compose.yml.template | sed "s/<ZOOKEEPER_PORT>/${ZOOKEEPER_PORT}/g" | sed "s/<KAFKA_PORT>/${KAFKA_PORT}/g" > "$DIR"/docker-compose.yml

docker-compose up -d
