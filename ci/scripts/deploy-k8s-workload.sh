#!/bin/sh -ex

kubo_deployment_dir="$(cd "$(dirname "$0")/../.."; pwd)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$kubo_deployment_dir/bosh.log"
export DEBUG=1

cp "$PWD/service-creds/service-ci-service-creds.yml" "$kubo_deployment_dir/ci/environments/gcp/"

$kubo_deployment_dir/bin/set_kubeconfig "$kubo_deployment_dir/ci/environments/gcp" ci-service
kubectl create -f $kubo_deployment_dir/ci/specs/nginx.yml
kubectl get pods --namespace=kube-system
