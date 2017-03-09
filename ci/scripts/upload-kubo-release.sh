#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"

cp -r "$PWD/git-kubo-release" "$PWD/kubo-release"

creds_path="${PWD}/s3-bosh-creds/creds.yml"
export BOSH_CLIENT="bosh_admin"
export BOSH_CLIENT_SECRET="$(bosh-cli int "$creds_path" --path /bosh_admin_client_secret)"
export BOSH_ENVIRONMENT="$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path /internal_ip)"
export BOSH_CA_CERT="$(bosh-cli int "${creds_path}" --path=/director_ssl/ca)"

cd $PWD/kubo-release
bosh-cli create-release --force --name "kubo-release" --tarball="$PWD/../s3-kubo-release-tarball/kubo-release-latest.tgz"
