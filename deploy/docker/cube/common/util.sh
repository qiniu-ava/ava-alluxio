#!/bin/bash

export Black="\033[0;30m"
export Red="\033[0;31m"
export Green="\033[0;32m"
export Orange="\033[0;33m"
export Blue="\033[0;34m"
export Cyan="\033[0;36m"
export Gray="\033[0;37m"
export NC="\033[0m"

# Bold
export BBlack='\033[1;30m'       # Black
export BRed='\033[1;31m'         # Red
export BGreen='\033[1;32m'       # Green
export BYellow='\033[1;33m'      # Yellow
export BBlue='\033[1;34m'        # Blue
export BCyan='\033[1;36m'        # Cyan
export BWhite='\033[1;37m'       # White


# max read workers in same node
export MAX_READ_WORKER_PER_NODE=4

join(){
    # If no arguments, do nothing.
    # This avoids confusing errors in some shells.
    if [ $# -eq 0 ]; then
        return
    fi

    local joiner="$1"
    shift

    while [ $# -gt 1 ]; do
        printf "%s%s" "$1" "$joiner"
        shift
    done

    printf '%s\n' "$1"
}

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
    path=/mnt/ramdisk-write-"$group"
  else
    path=/mnt/ramdisk-"$group"
  fi

  {
    if ! ([ -d "$path" ] && mountpoint "$path" ); then
      rm -rf "$path"
      mkdir -p "$path"
      mount -t ramfs -o size="$size" ramfs "$path"
      chmod a+w "$path"
      mkdir -p "$path"/data
    fi
  } > /dev/null
}

getAvailableSSD() {
  if isMacOS; then
    if [ $MACOS_VOLUME_PREFIX ]; then
      echo "${MACOS_VOLUME_PREFIX}"
    else
      echo -e ${BRed}"[ERROR]${NC} MACOS_VOLUME_PREFIX not set"
      exit 1
    fi
  fi

  # here we firmly assume that all the block device are mounted at /disk[0-9]+
  df -h | grep -E '/disk[0-9]+$' | sort -k 6,6 | awk '{print $1,$6}' | while read -r dev mountpoint; do
    if [ $(lsblk -o ROTA $dev | grep -v ROTA | awk '{print $1}') == "0"  ]; then
      echo $mountpoint
    fi
  done
}

getSSDAvailableQuota() {
  if [ $# -ne 1 ]; then
    echo -e "${BRed}please specify a SSD mountpoint${NC}"
    return 1
  fi

  origin_size=$(df -h $1 | grep "$1" | awk '{print $2}')
  {
    if ! which bc; then
      apt update
      apt install -y bc
    fi

    if echo $origin_size | grep T ; then
      origin_size=$(echo $origin_size | grep -oE '[0-9\.]+' | awk '{print $1"*10"}' | bc | awk -F'.' '{print $1"00"}')
    else
      origin_size=$(echo $origin_size | grep -oE '[0-9\.]+' | awk '{print $1"/100"}' | bc | awk -F'.' '{print $1"00"}')
    fi
  } &> /dev/null

  if [ "$1" = "/disk1" ]; then
    echo $(echo "${origin_size}*0.8/(2*${MAX_READ_WORKER_PER_NODE})" | bc | awk -F'.' '{print $1}')
  else
    echo $(echo "${origin_size}*0.8/${MAX_READ_WORKER_PER_NODE}" | bc | awk -F'.' '{print $1}')
  fi
}

# input:
#   video
#   /disk1
#   /disk2
#   /disk3
# output:
#   -v /disk1/alluxio/data-video/cachedisk:/opt/cachedisk1 -v /disk2/alluxio/data-video/cachedisk:/opt/cachedisk2 -v /disk3/alluxio/data-video/cachedisk:/opt/cachedisk3 
gen_volume_str_from_ssd_list() {
  if [ "$#" -le 1 ]; then
    echo "${BRed}[ERROR]${NC} arguments error, please give group_name and ssd_list"
    exit 1
  fi

  group=$1
  shift
  ss=""
  for disk in $@; do
    num=$(echo ${disk} | grep -oE "[0-9]+")
    ss="${ss} -v /disk${num}/alluxio/data-${group}/cachedisk:/opt/cachedisk${num} "
  done
  echo "${ss}"
}

# input:
#   /disk1
#   /disk2
#   /disk3
# output:
#   /opt/cachedisk1,/opt/cachedisk2,/opt/cachedisk3 
gen_path_str_from_ssd_list() {
  if [ "$#" -le 0 ]; then
    echo "${BRed}[ERROR]${NC} arguments error, please give ssd_list"
    exit 1
  fi

  echo $(join , $(echo $@ | sed 's/disk/opt\/cachedisk/g'))
}

# input:
#   /disk1
#   /disk2
#   /disk3
# example output:
#   200GB,200GB,200GB
gen_quota_str_from_ssd_list() {
  if [ "$#" -le 0 ]; then
    echo "${BRed}[ERROR]${NC} arguments error, please give ssd_list"
    exit 1
  fi
  echo $(join , $(for disk in "$@"; do echo $(getSSDAvailableQuota "${disk}")GB; done))
}

# only apply for unit in G
get_container_mem_size_from_jvm_size() {
  if [ $# -eq 1 ]; then
    if [ $(echo $1 | grep -cE "^\s*[0-9]+(G|g)i{0,1}\s*$") -eq 1 ]; then
      n=$(echo $1 | grep -oE "[0-9]+")
      echo $(printf "%.0f\n" $(echo "${n}*1.25+0.5" | bc))"G"
    else
      echo "invalid arguments $@ to get container mem size"
      exit 2
    fi
  else
    echo "invalid arguments $@ to get container mem size"
    exit 1
  fi
}
