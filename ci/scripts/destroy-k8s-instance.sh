#!/bin/sh -ex

kubo_deployment_dir="${PWD}/git-kubo-deployment"

creds_path="${PWD}/s3-bosh-creds/creds.yml"

export BOSH_CLIENT="bosh_admin"
export BOSH_CLIENT_SECRET="$(bosh-cli int "$creds_path" --path /bosh_admin_client_secret)"
export BOSH_ENVIRONMENT=$(bosh-cli int "${kubo_deployment_dir}/ci/environments/gcp/director.yml" --path /internal_ip)
export BOSH_CA_CERT="$(bosh-cli int "${creds_path}" --path=/director_ssl/ca)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${kubo_deployment_dir}/bosh.log"

bosh-cli -d ci-service -n delete-deployment
