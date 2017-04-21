#!/bin/bash -ex

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

cp "$PWD/s3-service-creds/ci-service-creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n):8844" --skip-tls-validation


export nginx_port=$(expr $(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port") + 1000)
export nginx_name="nginx-$(basename "${KUBO_ENVIRONMENT_DIR}")"

"${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service
kubectl create -f "${KUBO_DEPLOYMENT_DIR}/ci/specs/nginx.yml"
kubectl label services nginx http-route-sync=${nginx_name}
kubectl label services nginx tcp-route-sync=${nginx_port}
# wait for deployment to finish
kubectl rollout status deployment/nginx -w


check_tcp_route() {
  until wget -O - "http://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name"):${nginx_port}"; do
    sleep 1
  done
}
export -f check_tcp_route

check_http_route() {
  until wget -O - "http://${nginx_name}.$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")"; do
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
