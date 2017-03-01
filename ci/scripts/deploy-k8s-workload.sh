#!/bin/sh -ex

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"
export DEBUG=1

cp "$PWD/s3-service-creds/service-ci-service-creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"
cp "$PWD/s3-bosh-creds/creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"

credhub login -u credhub-user -p \
  "$(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "$kubo_deployment_dir/ci/environments/gcp/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation

$kubo_deployment_dir/bin/set_kubeconfig "$kubo_deployment_dir/ci/environments/gcp" ci-service
kubectl create -f $kubo_deployment_dir/ci/specs/nginx.yml
# wait for deployment to finish
kubectl rollout status deployment/nginx -w
