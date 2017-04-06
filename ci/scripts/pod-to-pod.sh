#!/bin/bash -ex

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${KUBO_DEPLOYMENT_DIR}/bosh.log"
export DEBUG=1

cp "$PWD/s3-service-creds/ci-service-creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

bosh_ca_cert=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip" | xargs echo -n)
credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password" | xargs echo -n)" \
  -s "https://${director_ip}:8844" --skip-tls-validation

"${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service
kubectl apply -f "${KUBO_DEPLOYMENT_DIR}/ci/specs/guestbook.yml"
# wait for deployment to finish
kubectl rollout status deployment/frontend -w
kubectl rollout status deployment/redis-master -w
kubectl rollout status deployment/redis-slave -w

worker_ip=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET=${client_secret} BOSH_CA_CERT="${bosh_ca_cert}" bosh-cli -e ${director_ip} vms | grep worker | head -n1 | awk '{print $4}' | xargs echo -n)
testvalue="hellothere$(date +'%N')"
wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=set&key=messages&value=${testvalue}"

if timeout 120 /bin/bash <<EOF
  until wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=get&key=messages" | grep ${testvalue}; do
    sleep 2
  done
EOF
then
  echo "Guestbook app is up and running"
else
  echo "Expected the sample guest book to display the stored value"
  exit 1
fi