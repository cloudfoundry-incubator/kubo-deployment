#!/bin/sh -ex

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"
export DEBUG=1

cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
cp -r "$PWD/git-kubo-release" "$PWD/kubo-release"

credhub login -u credhub-user -p \
  "$(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation
credhub set -n \
  "$(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/director.yml" --path="/director_name" | xargs echo -n)/ci-service/routing-cf-client-secret" \
  -t password -c $(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/director.yml" --path="/cf-tcp-router-name" | xargs echo -n) \
  -v "${ROUTING_CF_CLIENT_SECRET}" -O


"$kubo_deployment_dir/bin/set_bosh_alias" "$kubo_deployment_dir/ci/environments/gcp"
# Deploy k8s
"$kubo_deployment_dir/bin/deploy_k8s" "$kubo_deployment_dir/ci/environments/gcp" ci-service dev 
