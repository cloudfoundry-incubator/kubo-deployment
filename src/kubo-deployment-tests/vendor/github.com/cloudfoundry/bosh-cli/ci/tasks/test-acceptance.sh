#!/usr/bin/env bash

set -e -x

ensure_not_replace_value() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}

set +x
if [[ `whoami` != "root" ]]; then
  echo "acceptance tests must be run as a privileged user"
  exit 1
fi
set -x

export PATH=/usr/local/ruby/bin:/usr/local/go/bin:$PATH
export GOPATH=$PWD/gopath

export BOSH_INIT_CPI_RELEASE_PATH=`ls $PWD/cpi-release/*.tgz`
export BOSH_INIT_CPI_RELEASE_URL=""
export BOSH_INIT_CPI_RELEASE_SHA1=""

cd $GOPATH/src/github.com/cloudfoundry/bosh-cli

source /etc/profile.d/chruby.sh
chruby 2.3.1

start-garden

base=$PWD ./bin/test-acceptance-with-garden
