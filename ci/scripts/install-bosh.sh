#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

printenv GCP_SERVICE_ACCOUNT > key.json
set -x
gcloud auth activate-service-account bosh-deployer@cf-pcf-kubo.iam.gserviceaccount.com --key-file=key.json
"$kubo_deployment_dir/bin/deploy_bosh" "$kubo_deployment_dir/ci/environments/gcp"
cp "$kubo_deployment_dir/ci/environments/gcp/creds.yml" bosh-creds/
cp "$kubo_deployment_dir/ci/environments/gcp/*.json" bosh-state/
