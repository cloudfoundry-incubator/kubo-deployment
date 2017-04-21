#!/bin/sh -e

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

tarball_name=$(ls $PWD/s3-kubo-release-tarball/kubo-release*.tgz | head -n1)

cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

cp "$tarball_name" "${KUBO_DEPLOYMENT_DIR}/../kubo-release.tgz"

"${KUBO_DEPLOYMENT_DIR}/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
"${KUBO_DEPLOYMENT_DIR}/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" ci-service local

cp "${KUBO_ENVIRONMENT_DIR}/ci-service-creds.yml" "$PWD/service-creds/ci-service-creds.yml"
