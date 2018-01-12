#!/usr/bin/env bash

set -euo pipefail

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

get_tcp_hostname() {
	# Looks up the TCP hostname using `cf domains`
	# See https://docs.cloudfoundry.org/adminguide/enabling-tcp-routing.html#-configure-cf-with-your-tcp-domain
	# to
	echo "$(cf domains | grep 'tcp$' | awk '{print $1}')"
}

get_routing_client_secret() {
 	credhub get -n "$BOSH_ENVIRONMENT/cf/uaa_clients_routing_api_client_secret" --output-json | jq -r .value
}

get_nats_password() {
	credhub get -n "$BOSH_ENVIRONMENT/cf/nats_password" --output-json | jq -r .value
}

get_nats_ips_json() {
	bosh instances -d "$CF_DEPLOYMENT" --json | jq '[.Tables[].Rows[] | select(.instance | contains("nats")) | .ips]'
}

main() {
	local cf_manifest
	local system_domain
	local app_domain

	cf_manifest=$(bosh manifest -d "$CF_DEPLOYMENT")
	system_domain=$(echo "$cf_manifest" | bosh int - --path /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/system_domain)
	app_domain=$(echo "$cf_manifest" | bosh int - --path /instance_groups/name=api/jobs/name=cloud_controller_ng/properties/app_domains/0)

	cat <<YAML
kubernetes_master_host: $(get_tcp_hostname)
kubernetes_master_port: 8443
routing_cf_api_url: https://api.$system_domain
routing_cf_uaa_url: https://uaa.$system_domain
routing_cf_app_domain_name: $app_domain
routing_cf_client_id: routing_api_client
routing_cf_client_secret: "$(get_routing_client_secret)"
routing_cf_nats_port: 4222
routing_cf_nats_username: nats
routing_cf_nats_password: "$(get_nats_password)"
routing_cf_nats_internal_ips: $(get_nats_ips_json)
YAML

}

[[ "$0" == "${BASH_SOURCE[0]}" ]] && main "$@"
