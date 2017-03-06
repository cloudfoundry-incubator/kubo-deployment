#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/environment.sh"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"

cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
cp  "$PWD/s3-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"

# Destroy KuBOSH
"${KUBO_DEPLOYMENT_DIR}/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
