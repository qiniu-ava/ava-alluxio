#!/bin/bash

CONFIGURE_LOCK="/opt/alluxio/configure.lock"

ERROR_BAD_COMMAND=1000
ERROR_COMMAND_NOT_FOUND=1001

ERROR_UNKNOW_MOUNT_ERROR=1100
ERROR_FAILED_CREATE_PREFIX=1101
ERROR_BAD_UFS_TYPE=1102

log_dir=/var/log/alluxio

if [ -z "$ENDPOINT" ]; then
  ENDPOINT="iovip.qbox.me"
fi

# shellcheck disable=SC2164
CUR_DIR=$(cd "$( dirname "$( readlink "$0" || echo "$0" )" )"; pwd)
USAGE="${CUR_DIR}/alluxio-flex-volume.sh <command> <options>
  command:
    mount [--mode=<read_write_mode>] --group=<group_name> --uid=<uid> --ak=<access_key> --sk=<secret_key> --bucket=<bucket_name> [--prefix=<key_prefix>] --domain=<bucket_domain> --local_path=<local_path>
    umount <local_path>
"

cluster_map=$(cat <<-END
{
  "jw": {
    "name": "juewa",
    "ip_prefix": "10.18.21",
    "ip_exclude": ""
  },
  "juewa": {
    "ip_prefix": "10.18.21",
    "ip_exclude": ""
  },
  "jq": {
    "name": "jq",
    "ip_prefix": "192.168.",
    "ip_exclude": "192.168.212"
  }
}
END
)

get_cluster_name() {
  cluster_short_cut=$(hostname | sed 's/gpu//g' | sed 's/[0-9]\+//g')
  cluster_name=$(echo $cluster_map | jq --arg sc $cluster_short_cut '.$sc')

  if [ "$cluster_name" = "" ]; then
    cluster_name=jq
  fi

  echo $cluster_name
}

print_usage() {
  echo -e "${USAGE}" >&2
  if [ "$?" = "" ]; then
    exit ERROR_BAD_COMMAND
  else
    exit $?
  fi
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
  cluster=$(get_cluster_name)
  ip_prefix=$(echo "$cluster_map" | jq --arg cluster $cluster '.$cluster.ip_prefix')
  ip_exclude=$(echo "$cluster_map" | jq --arg cluster $cluster '.$cluster.ip_exclude')

  if is_linux; then
    if [ "$ip_exclude" ]; then
      ifconfig | grep -E "^\\s+inet " | grep -- "${ip_prefix}" | grep -v "${ip_exclude}" | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    else
      ifconfig | grep -E "^\\s+inet " | grep -- "${ip_prefix}" | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    fi
    # if [[ $(ifconfig | grep -E "^\\s+inet " | grep -c "192.168.") -lt 2 ]]; then
    #   ifconfig | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    # else
    #   ifconfig | grep -E "^\\s+inet " | grep "192.168." | grep -v "192.168.212" | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
    # fi
  elif is_macos; then
    ifconfig | grep -E -A6 "^en0:" | grep -E "^\\s+inet " | grep -o -E  '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}' | head -n 1
  else
    echo "not compatible OS"
    return 1
  fi
}

setup_log_dir() {
  local bucket
  local ufs_type

  for i in "$@"; do
    case $i in
      --bucket=*)
        bucket="${i#*=}"
      ;;
      --ufs_type=*)
        ufs_type="${i#*=}"
        ;;
    esac
  done
  log_dir="/var/log/alluxio/${ufs_type}-${bucket}-$(date +%s)"

  mkdir -p $log_dir
}

acquire_configure() {
  while [ -f "${CONFIGURE_LOCK}" ]; do
    sleep 1
  done

  cp "$CUR_DIR/conf/alluxio-env-$1.sh" "$CUR_DIR/conf/alluxio-env.sh"
  cp "$CUR_DIR/conf/alluxio-site-$1.properties" "$CUR_DIR/conf/alluxio-site.properties"
  echo "alluxio.logs.dir=${log_dir}" >> $CUR_DIR/conf/alluxio-site.properties
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
  local ufs_type
  local alluxio_uri_prefix

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
      --ufs_type=*)
        ufs_type="${i#*=}"
        ;;
      *)
        echo "[WARNING] unknow argument ${i}"
      ;;
    esac
  done

  cluster=$(get_cluster_name)
  if [[ -f "$CUR_DIR/conf/alluxio-env-$group.sh" && -f "$CUR_DIR/conf/alluxio-site-$group.properties" ]]; then
    group="$group"
  elif [ "$cluster" != "jq" ]; then
    group="$cluster-default"
  else
    group=default
  fi

  if [ "$cluster" != "jq" ]; then
    alluxio_uri_prefix="/${ufs_type}"
  else
    if [ "$group" = "default" ]; then
      alluxio_uri_prefix="/ava/qn-bucket/$uid"
    else
      alluxio_uri_prefix="/$uid"
    fi
  fi

  if [ "$mode" = "" ]; then
    mode=ro
  fi

  # create user path in alluxio if it doesn't exist
  acquire_configure "$group"
  if ! "${CUR_DIR}"/bin/alluxio fs stat "$alluxio_uri_prefix" ; then
    if ! "${CUR_DIR}"/bin/alluxio fs mkdir "$alluxio_uri_prefix" ; then
      echo "failed to create path for $uid in $group"
      release_configure
      exit $ERROR_FAILED_CREATE_PREFIX
    fi
  fi
  release_configure

  # mount bucket in alluxio
  acquire_configure "$group"
  if [ "$cluster" != "jq" ]; then
    if [ "$ufs_type" = "kodo" ]; then
      "${CUR_DIR}"/bin/alluxio fs mount \
        --option fs.kodo.accesskey="${ak}" \
        --option fs.kodo.secretkey="${sk}" \
        --option fs.kodo.endpoint="${ENDPOINT}" \
        --option fs.kodo.downloadhost="${domain}" \
        "$alluxio_uri_prefix/$bucket" \
        "kodo://$bucket"
    elif [ "$ufs_type" = "oss" ]; then
      "${CUR_DIR}"/bin/alluxio fs mount \
        --option fs.oss.accessKeyId="${ak}" \
        --option fs.oss.accessKeySecret="${sk}" \
        --option fs.oss.endpoint="${domain}" \
        --option fs.oss.userId="${uid}" \
        "$alluxio_uri_prefix/$bucket" \
        "oss://$bucket"
    else
      release_configure
      echo "unsupported under filesystem type: $ufs_type"
      exit $ERROR_BAD_UFS_TYPE
    fi
  else
    "${CUR_DIR}"/bin/alluxio fs mount \
      --option fs.oss.accessKeyId="${ak}" \
      --option fs.oss.accessKeySecret="${sk}" \
      --option fs.oss.endpoint="${domain}" \
      --option fs.oss.userId="${uid}" \
      "$alluxio_uri_prefix/$bucket" \
      "oss://$bucket"
  fi
  release_configure

  # mount alluxio path to local path
  acquire_configure "$group"
  if [ -n "$prefix" ]; then
    "${CUR_DIR}"/integration/fuse/bin/alluxio-fuse mount "$local_path" "$alluxio_uri_prefix/$bucket" -o "$mode"
  else
    "${CUR_DIR}"/integration/fuse/bin/alluxio-fuse mount "$local_path" "$alluxio_uri_prefix/$bucket/$prefix" -o "$mode"
  fi
  release_configure
}

flex_volume_unmount() {
  # unmount local path
  if [ "$#" -ne 1 ]; then
    print_usage $ERROR_BAD_COMMAND
  fi

  if mount | grep "$1" ; then
    "${CUR_DIR}"/integration/fuse/bin/alluxio-fuse umount "$1"
  else
    echo "$1 is not a mountpoint"
  fi
}

if [ "$#" -lt 1 ]; then
  print_usage $ERROR_BAD_COMMAND
fi

case "$1" in
  mount)
    shift
    setup_log_dir "$@"
    flex_volume_mount "$@"
  ;;
  umount)
    shift
    flex_volume_unmount "$@"
  ;;
  *)
    print_usage $ERROR_COMMAND_NOT_FOUND
  ;;
esac
