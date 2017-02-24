#!/bin/sh -ex

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"
export DEBUG=1

cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
cp  "$PWD/s3-bosh-state/state.json" "$kubo_deployment_dir/ci/environments/gcp/"
mv "$PWD/git-kubo-release" "$PWD/kubo-release"

"$kubo_deployment_dir/bin/set_bosh_alias" "$kubo_deployment_dir/ci/environments/gcp"
# Deploy k8s
"$kubo_deployment_dir/bin/deploy_k8s" "$kubo_deployment_dir/ci/environments/gcp" ci-service dev 
