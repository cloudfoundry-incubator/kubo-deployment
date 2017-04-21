#!/bin/sh -ex

creds_path="${PWD}/s3-bosh-creds/creds.yml"
. "$(dirname "$0")/lib/environment.sh"

export BOSH_CLIENT="bosh_admin"
export BOSH_CLIENT_SECRET="$(bosh-cli int "$creds_path" --path /bosh_admin_client_secret)"
export BOSH_ENVIRONMENT="$(bosh-cli int "kubo-lock/metadata" --path /internal_ip)"
export BOSH_CA_CERT="$(bosh-cli int "${creds_path}" --path=/director_ssl/ca)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"

bosh-cli -d ci-service -n delete-deployment
