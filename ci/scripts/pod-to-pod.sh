#!/bin/sh -ex

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

cp "$PWD/s3-service-creds/service-ci-service-creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n)
credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://${director_ip}:8844" --skip-tls-validation

"${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service
kubectl create -f "${KUBO_DEPLOYMENT_DIR}/ci/specs/guestbook.yml"
# wait for deployment to finish
kubectl rollout status deployment/nginx -w

worker_ip=$(bosh-cli -e ${director_ip} vms | grep worker | head -n1 | awk '{print $4}' | xargs echo -n)
curl http://${worker_ip}:30303/guestbook.php?cmd=set&key=messages&value=hellothere
result=$(curl http://${worker_ip}:30303 | grep hellothere)
if [ -z result ] ; then
  echo "Expected the sample guestbook to display the stored value"
  exit 1
fi