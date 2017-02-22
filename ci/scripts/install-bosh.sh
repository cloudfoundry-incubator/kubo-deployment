#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"

# Deploy Bosh++
"$kubo_deployment_dir/bin/deploy_bosh" "$kubo_deployment_dir/ci/environments/gcp" "$PWD/key.json"

cp "$kubo_deployment_dir/ci/environments/gcp/creds.yml" "$PWD/bosh-creds/"
cp "$kubo_deployment_dir/ci/environments/gcp/state.json" "$PWD/bosh-state/"

