#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"

cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
cp  "$PWD/s3-bosh-state/state.json" "$kubo_deployment_dir/ci/environments/gcp/"

# Destroy Bosh++
"$kubo_deployment_dir/bin/destroy_bosh" "$kubo_deployment_dir/ci/environments/gcp" "$PWD/key.json"
