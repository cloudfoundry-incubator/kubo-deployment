#!/bin/bash

set -eu

# Requires access to Credhub via `credhub` CLI to lookup secrets
# Usage:
#
#   Either run once to get a copy:
#     ./kubo-deployment/manifests/helper/cf-routing-vars.sh > cf-vars.yml
#     bosh deploy kubo-deployment/manifests/cfcr.yml \
#       -o kubo-deployment/manifests/ops-files/cf-routing.yml \
#       -l cf-vars.yml
#
#  Run using process substitution:
#     bosh deploy kubo-deployment/manifests/cfcr.yml \
#       -o kubo-deployment/manifests/ops-files/cf-routing.yml \
#       -l <(./kubo-deployment/manifests/helper/cf-routing-vars.sh)
#

: ${BOSH_ENVIRONMENT:?required}
CF_DEPLOYMENT=${CF_DEPLOYMENT:-cf}
cf_manifest=$(bosh manifest -d $CF_DEPLOYMENT)

# Looks up the TCP hostname using `cf domains`
# See https://docs.cloudfoundry.org/adminguide/enabling-tcp-routing.html#-configure-cf-with-your-tcp-domain
# to
tcp_hostname=${tcp_hostname:-$(cf domains | grep "tcp$" | awk '{print $1}')}

system_domain=$(echo "$cf_manifest" | bosh int - --path /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/system_domain)
app_domain=$(echo "$cf_manifest" | bosh int - --path /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/app_domains/0)
routing_cf_client_secret=${routing_cf_client_secret:-$(credhub get -n $BOSH_ENVIRONMENT/cf/uaa_clients_routing_api_client_secret --output-json | jq -r .value)}
nats_password=${nats_password:-$(credhub get -n $BOSH_ENVIRONMENT/cf/nats_password --output-json | jq -r .value)}
nats_ips=($(bosh instances -d $CF_DEPLOYMENT | grep nats | awk '{print $4}'))
nats_ips_list=""
for nats_ip in "$nats_ips"; do
  if [[ "${nats_ips_list:-empty}" != "empty" ]]; then
    nats_ips_list="${nats_ips_list},"
  fi
  nats_ips_list="${nats_ips_list}${nats_ip}"
done
nats_ip_json="[$nats_ips_list]"

cat <<YAML
kubernetes_master_host: ${tcp_hostname}
kubernetes_master_port: 8443
routing_cf_api_url: https://api.$system_domain
routing_cf_uaa_url: https://uaa.$system_domain
routing_cf_app_domain_name: $app_domain
routing_cf_client_id: routing_api_client
routing-cf-client-secret: "${routing_cf_client_secret}"
routing_cf_nats_internal_ips: ${nats_ip_json}
routing_cf_nats_port: 4222
routing_cf_nats_username: nats
routing-cf-nats-password: "${nats_password}"
YAML
