#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
version=$(cat kubo-version/version)

cd git-kubo-release
bosh-cli create-release --name "kubo" --tarball="../kubo-release/kubo-release-latest.tgz" --version=${version}
