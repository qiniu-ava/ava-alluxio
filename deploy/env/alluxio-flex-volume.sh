#!/bin/bash

bold=$(tput bold)
normal=$(tput sgr0)

CONFIGURE_LOCK="/opt/alluxio/configure.lock"

CUR_DIR=$(cd "$( dirname "$( readlink "$0" || echo "$0" )" )"; pwd)
USAGE="${CUR_DIR}/alluxio-flex-volume.sh <command> <options>
  command:
    ${bold}mount${normal} [--group=<group_name>] [--mode=<read_write_mode>] --uid=<uid> --ak=<access_key> --sk=<secret_key> --bucket=<bucket_name> --domain=<bucket_domain> --local_path=<local_path>
    ${bold}unmount${normal} <local_path>
"

print_usage() {
  echo -e "${USAGE}" >&2
  exit 1
}

acquire_configure() {
  while [ -f "${CONFIGURE_LOCK}" ]; do
    sleep 1
  done

  local group
  group=$1

  if [ -f "$CUR_DIR/conf/alluxio-env-$group.sh" && -f "$CUR_DIR/conf/alluxio-site-$group.properties" ]; then
    cp "$CUR_DIR/conf/alluxio-env-$group.sh" "$CUR_DIR/conf/alluxio-env.sh"
    cp "$CUR_DIR/conf/alluxio-site-$group.properties" "$CUR_DIR/conf/alluxio-site.properties"
  else
    cp "$CUR_DIR/conf/alluxio-env-default.sh" "$CUR_DIR/conf/alluxio-env.sh"
    cp "$CUR_DIR/conf/alluxio-site-default.properties" "$CUR_DIR/conf/alluxio-site.properties"
  fi

  if [ $? -ne 0 ]; then
    echo "failed to copy alluxio configure"
    exit 10
  fi

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
  local domain
  local mode
  local local_path
  local alluxio_uid_path

  for i in $@; do
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
    esac
  done

  if [ "$group" = "" ]; then
    group=default
    alluxio_uid_path="/ava/qn-bucket/$uid"
  else
    alluxio_uid_path="/$uid"
  fi

  if [ "$mode" = "" ]; then
    mode=ro
  fi

  # create path in alluxio
  acquire_configure "$group"
  ${CUR_DIR}/bin/alluxio fs mkdir "$alluxio_uid_path"
  ${CUR_DIR}/bin/alluxio fs stat "$alluxio_uid_path"
  if [ $? -ne 0 ]; then
    echo "failed to create path for $uid in $group"
    release_configure
    exit 20
  fi
  release_configure

  # mount bucket in alluxio
  acquire_configure "$group"
  ${CUR_DIR}/bin/alluxio fs mount \
    --option fs.oss.accessKeyId=${ak} \
    --option fs.oss.accessKeySecret=${sk} \
    --option fs.oss.endpoint=${domain} \
    --option fs.oss.userId=${uid} \
    "$alluxio_uid_path/$bucket" \
    "oss://$bucket"
  release_configure

  # mount alluxio path to local path
  acquire_configure "$group"
  ${CUR_DIR}/integration/fuse/bin/alluxio-fuse mount "$local_path" "$alluxio_uid_path/$bucket" -o "$mode"
  release_configure
}

flex_volume_unmount() {
  # unmount local path
  if [ "$#" -ne 1 ]; then
    print_usage
  fi

  ${CUR_DIR}/integration/fuse/bin/alluxio-fuse umount "$1"
}

if [ "$#" -lt 1 ]; then
  print_usage
fi

case "$1" in
  mount)
    shift
    flex_volume_mount $@
  ;;
  unmount)
    shift
    flex_volume_unmount $@
  ;;
  *)
    print_usage
  ;;
esac
