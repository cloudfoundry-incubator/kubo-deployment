#!/bin/sh -e

. "$(dirname "$0")/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp -r "$PWD/git-kubo-release" "$PWD/kubo-release"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation
credhub set -n \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name" | xargs echo -n)/ci-service/routing-cf-client-secret" \
  -t password -c "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name" | xargs echo -n)" \
  -v "${ROUTING_CF_CLIENT_SECRET}" -O > /dev/null


"${KUBO_DEPLOYMENT_DIR}/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
# Deploy k8s
"${KUBO_DEPLOYMENT_DIR}/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" ci-service dev

cp "${KUBO_ENVIRONMENT_DIR}/service-ci-service-creds.yml" "$PWD/service-creds/service-ci-service-creds.yml"
