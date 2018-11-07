#!/bin/bash

isLinux() {
  test $(uname -s) = "Linux"
  return $?
}

isMacOS() {
  test $(uname -s) = "Darwin"
  return $?
}

# only apply for Qiniu office WIFI and jq cluster network settings
getMyIP() {
  if isLinux; then
    ifconfig | grep -E "^\\s+inet " | grep "192.168." | grep -v "192.168.212" | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
  elif isMacOS; then
    ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
  else
    echo "not compatible OS"
    return 1
  fi
}

ensureGroupRamdisk() {
  group=$1
  size=$2
  mode=$3
  local path

  if [ "$mode" == "write" ]; then
    path=/mnt/ramdisk-writer-"$group"
  else
    path=/mnt/ramdisk-"$group"
  fi

  if ! ([ -d "$path" ] && mountpoint "$path" ); then
    rm -rf "$path"
    mkdir -p "$path"
    mount -t ramfs -o size="$size" ramfs "$path"
    chmod a+w "$path"
    mkdir -p "$path"/data
  fi
}
