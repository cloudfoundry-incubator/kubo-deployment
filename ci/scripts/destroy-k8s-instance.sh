#!/bin/sh -ex

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"
#cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
$kubo_deployment_dir/ci/scripts/director-environment.sh "$kubo_deployment_dir/ci/environments/gcp"
echo "$kubo_deployment_dir/ci/environments/gcp"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"
export DEBUG=1

bosh-cli -e gcp -d ci-service deld
