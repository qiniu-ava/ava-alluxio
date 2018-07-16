#!/bin/bash

set -x

need_installation=""

if [[ $(command -v mdl) == "" ]];then
  echo "plz install markdownlint with `gem install mdl`"
  need_installation=true
fi

if [[ $(command -v shellcheck) == "" ]];then
  echo "plz install shellcheck with `brew install shellcheck` on macos"
  echo "or `apt-get install shellcheck` on ubuntu"
  need_installation=true
fi

if [ "$need_installation" == "true" ]; then
  exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# lint markdown
cd $DIR/../..
mdl ./*.md ./docs
cd $DIR

# lint shell
cd $DIR/../..
shellcheck `find . -type d \( -path ./.tmp -o -path ./alluxio -o -path ./kodo -o -path ./tools/golang/src/qiniu.com/vendor \) -prune -o -name "*.sh" -print`
cd $DIR
