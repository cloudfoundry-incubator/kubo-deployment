#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/lib/environment.sh"

set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"

bosh-cli int kubo-lock/metadata --path=/gcp_service_account > "$PWD/key.json"

cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

"${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"

cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
