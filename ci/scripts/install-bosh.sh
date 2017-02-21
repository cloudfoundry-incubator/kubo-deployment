#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
gcloud auth activate-service-account bosh-deployer@cf-pcf-kubo.iam.gserviceaccount.com --key-file="$PWD/key.json"
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"

# Deploy Bosh++
"$kubo_deployment_dir/bin/deploy_bosh" "$kubo_deployment_dir/ci/environments/gcp" "$PWD/key.json"

# TODO: Stash these creds so the destroy step can be seperate
echo "TODO" > bosh-creds/creds.yml
echo "TODO" > bosh-state/state.json

# Destroy Bosh++
"$kubo_deployment_dir/bin/destroy_bosh" "$kubo_deployment_dir/ci/environments/gcp" "$PWD/key.json"
