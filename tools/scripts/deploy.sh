#!/bin/bash

RED='\033[0;31m'
BRED='\033[1;31m'
GREEN="\033[0;32m"
BGREEN='\033[1;32m'
NC='\033[0m' # No Color

version=` \
  grep 'declare-version' ./tools/golang/src/qiniu.com/app/avio/main.go \
    | awk -F'"' '{print $2}' \
`

check_version() {
  echo -e "did you login qrsctl as ${BRED}ava-test${NC} and are you sure to deploy avio in version ${BRED}${version}${NC} ? y/n"
  read yn
  case ${yn} in
    [Yy])
      exit 0
      ;;
    [Nn])
      exit 1
      ;;
  esac
}

deploy_avio() {
  echo -e "${GREEN}uploading executable binaries to kodo${NC}"
  qrsctl put devtools ava/cli/avio/${version}/avio-linux ./tools/golang/bin/linux_amd64/avio
  qrsctl put devtools ava/cli/avio/avio-linux ./tools/golang/bin/linux_amd64/avio
  qrsctl put devtools ava/cli/avio ./tools/golang/bin/linux_amd64/avio
  qrsctl put devtools ava/cli/avio/${version}/avio-darwin ./tools/golang/bin/avio
}

for i in "$@"; do
  case $i in
    --check-version)
      check_version
      exit $?
    ;;
  esac
done

deploy_avio
