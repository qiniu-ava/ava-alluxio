#!/bin/bash

CONFIGURE_LOCK="/opt/alluxio/configure.lock"

# shellcheck disable=SC2164
CUR_DIR=$(cd "$( dirname "$( readlink "$0" || echo "$0" )" )"; pwd)
USAGE="${CUR_DIR}/alluxio-flex-volume.sh <command> <options>
  command:
    mount [--mode=<read_write_mode>] --group=<group_name> --uid=<uid> --ak=<access_key> --sk=<secret_key> --bucket=<bucket_name> [--prefix=<key_prefix>] --domain=<bucket_domain> --local_path=<local_path>
    umount <local_path>
"

print_usage() {
  echo -e "${USAGE}" >&2
  exit 1
}

is_linux() {
  test $(uname -s) = "Linux"
  return $?
}

is_macos() {
  test $(uname -s) = "Darwin"
  return $?
}

# now we only consider ubuntu distribution
get_my_ip() {
  if is_linux; then
    if [[ $(ifconfig | grep -E "^\\s+inet " | grep -c "192.168.") -lt 2 ]]; then
      ifconfig | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    else
      ifconfig | grep -E "^\\s+inet " | grep "192.168." | grep -v "192.168.212" | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    fi
  elif is_macos; then
    ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
  else
    echo "not compatible OS"
    return 1
  fi
}

acquire_configure() {
  while [ -f "${CONFIGURE_LOCK}" ]; do
    sleep 1
  done

  cp "$CUR_DIR/conf/alluxio-env-$1.sh" "$CUR_DIR/conf/alluxio-env.sh"
  cp "$CUR_DIR/conf/alluxio-site-$1.properties" "$CUR_DIR/conf/alluxio-site.properties"
  echo "alluxio.locality.node=$(get_my_ip)" >> $CUR_DIR/conf/alluxio-site.properties
  echo "alluxio.user.hostname=$(get_my_ip)" >> $CUR_DIR/conf/alluxio-site.properties

  echo "this file is somelike a lock for alluxio configures" > "${CONFIGURE_LOCK}"
}

release_configure() {
  rm -rf "${CONFIGURE_LOCK}"
}

flex_volume_mount() {
  local group
  local uid
  local ak
  local sk
  local bucket
  local prefix
  local domain
  local mode
  local local_path
  local alluxio_uid_path

  for i in "$@"; do
    case $i in
      --group=*)
        group="${i#*=}"
      ;;
      --uid=*)
        uid="${i#*=}"
      ;;
      --ak=*)
        ak="${i#*=}"
      ;;
      --sk=*)
        sk="${i#*=}"
      ;;
      --bucket=*)
        bucket="${i#*=}"
      ;;
    --prefix=*)
      prefix="${i#*=}"
      ;;
      --domain=*)
        domain="${i#*=}"
      ;;
      --mode=*)
        mode="${i#*=}"
      ;;
      --local_path=*)
        local_path="${i#*=}"
      ;;
      *)
        echo "[WARNING] unknow argument ${i}"
      ;;
    esac
  done

  if [[ -f "$CUR_DIR/conf/alluxio-env-$group.sh" && -f "$CUR_DIR/conf/alluxio-site-$group.properties" ]]; then
    group="$group"
  else
    group=default
  fi

  if [ "$group" = "default" ]; then
    alluxio_uid_path="/ava/qn-bucket/$uid"
  else
    alluxio_uid_path="/$uid"
  fi

  if [ "$mode" = "" ]; then
    mode=ro
  fi

  # create user path in alluxio if it doesn't exist
  acquire_configure "$group"
  if ! "${CUR_DIR}"/bin/alluxio fs stat "$alluxio_uid_path" ; then
    if ! "${CUR_DIR}"/bin/alluxio fs mkdir "$alluxio_uid_path" ; then
      echo "failed to create path for $uid in $group"
      release_configure
      exit 10
    fi
  fi
  release_configure

  fullpath="$bucket"
  if [[ ! -z $prefix ]]; then
    fullpath="$fullpath/$prefix"
  fi

  # mount bucket in alluxio
  acquire_configure "$group"
  "${CUR_DIR}"/bin/alluxio fs mount \
    --option fs.oss.accessKeyId="${ak}" \
    --option fs.oss.accessKeySecret="${sk}" \
    --option fs.oss.endpoint="${domain}" \
    --option fs.oss.userId="${uid}" \
    "$alluxio_uid_path/$fullpath" \
    "kodo://$fullpath"
  release_configure

  # mount alluxio path to local path
  acquire_configure "$group"
  "${CUR_DIR}"/integration/fuse/bin/alluxio-fuse mount "$local_path" "$alluxio_uid_path/$fullpath" -o "$mode"
  release_configure
}

flex_volume_unmount() {
  # unmount local path
  if [ "$#" -ne 1 ]; then
    print_usage
  fi

  if mount | grep "$1" ; then
    "${CUR_DIR}"/integration/fuse/bin/alluxio-fuse umount "$1"
  else
    echo "$1 is not a mountpoint"
  fi
}

if [ "$#" -lt 1 ]; then
  print_usage
fi

case "$1" in
  mount)
    shift
    flex_volume_mount "$@"
  ;;
  umount)
    shift
    flex_volume_unmount "$@"
  ;;
  *)
    print_usage
  ;;
esac
