#!/bin/sh -e

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

version=$(cat kubo-version/version)

cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

cp "$PWD/s3-kubo-release-tarball/kubo-release-${version}.tgz" "${KUBO_DEPLOYMENT_DIR}/../kubo-release.tgz"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation
credhub set -n \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name" | xargs echo -n)/ci-service/routing-cf-client-secret" \
  -t password -c "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name" | xargs echo -n)" \
  -v "${ROUTING_CF_CLIENT_SECRET}" -O > /dev/null


"${KUBO_DEPLOYMENT_DIR}/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
# Deploy k8s
"${KUBO_DEPLOYMENT_DIR}/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" ci-service local

cp "${KUBO_ENVIRONMENT_DIR}/service-ci-service-creds.yml" "$PWD/service-creds/service-ci-service-creds.yml"
