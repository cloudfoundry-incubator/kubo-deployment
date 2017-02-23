#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"

cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
cp  "$PWD/s3-bosh-state/state.json" "$kubo_deployment_dir/ci/environments/gcp/"
cp -r "$PWD/git-kubo-release" "$PWD/kubo-release"

DEBUG=1 "$kubo_deployment_dir/bin/set_bosh_alias" "$kubo_deployment_dir/ci/environments/gcp"
# Deploy k8s
DEBUG=1 "$kubo_deployment_dir/bin/deploy_k8s" "$kubo_deployment_dir/ci/environments/gcp" ci-service dev 

cp "$kubo_deployment_dir/ci/environments/gcp/service-ci-service-creds.yml" "$PWD/s3-k8s-service-creds"
cp "$kubo_deployment_dir/ci/environments/gcp/state.json" "$PWD/s3-k8s-bosh-state"
