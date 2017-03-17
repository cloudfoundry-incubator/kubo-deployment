#!/bin/bash -ex

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

cp "$PWD/s3-service-creds/ci-service-creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation

"${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service
kubectl create -f "${KUBO_DEPLOYMENT_DIR}/ci/specs/nginx.yml"
# wait for deployment to finish
kubectl rollout status deployment/nginx -w


check_tcp_route() {
  until curl $(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name"):$(kubectl describe service nginx | grep NodePort | tr -dc '0-9'); do
    sleep 1
  done
}
export -f check_tcp_route

check_http_route() {
  until $(curl --output /dev/null --silent --head --fail nginx.default.$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")); do
    sleep 1
  done
}
export -f check_http_route

if timeout '60s' bash -c check_tcp_route
then
  echo 'TCP route sync is working'
else
  echo 'Nginx TCP route is not exposed :('
  exit 1
fi

if timeout '60s' bash -c check_http_route
then
  echo 'HTTP route sync is working'
else
  echo 'Nginx HTTP route is not exposed :('
  exit 1
fi
