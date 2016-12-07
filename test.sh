#!/bin/bash

set -eu

vars_store_prefix=/tmp/bosh-deployment-test

clean_tmp() {
  rm -f ${vars_store_prefix}.*
}

trap clean_tmp EXIT

echo "- AWS"
bosh interpolate bosh.yml \
  -o aws/cpi.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v access_key_id=test \
  -v secret_access_key=test \
  -v az=test \
  -v region=test \
  -v default_key_name=test \
  -v default_security_groups=[test] \
  -v private_key=test \
  -v subnet_id=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- AWS with UAA"
bosh interpolate bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v access_key_id=test \
  -v secret_access_key=test \
  -v az=test \
  -v region=test \
  -v default_key_name=test \
  -v default_security_groups=[test] \
  -v private_key=test \
  -v subnet_id=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- AWS with UAA for BOSH development"
bosh interpolate bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  -o bosh-dev.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_ip=test \
  -v access_key_id=test \
  -v secret_access_key=test \
  -v region=test \
  -v default_key_name=test \
  -v default_security_groups=[test] \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- AWS (cloud-config)"
bosh interpolate aws/cloud-config.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v az=test \
  -v subnet_id=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- GCP"
bosh interpolate bosh.yml \
  -o gcp/cpi.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- GCP with UAA"
bosh interpolate bosh.yml \
  -o gcp/cpi.yml \
  -o uaa.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- GCP with BOSH Lite"
bosh interpolate bosh.yml \
  -o gcp/cpi.yml \
  -o bosh-lite.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- GCP (cloud-config)"
bosh interpolate gcp/cloud-config.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- Openstack"
bosh interpolate bosh.yml \
  -o openstack/cpi.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v auth_url=test \
  -v az=test \
  -v default_key_name=test \
  -v default_security_groups=test \
  -v net_id=test \
  -v openstack_password=test \
  -v openstack_username=test \
  -v private_key=test \
  -v region=test \
  -v tenant=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- Openstack (cloud-config)"
bosh interpolate openstack/cloud-config.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v az=test \
  -v net_id=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- vSphere"
bosh interpolate bosh.yml \
  -o vsphere/cpi.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v network_name=test \
  -v vcenter_dc=test \
  -v vcenter_ds=test \
  -v vcenter_ip=test \
  -v vcenter_user=test \
  -v vcenter_password=test \
  -v vcenter_templates=test \
  -v vcenter_vms=test \
  -v vcenter_disks=test \
  -v vcenter_cluster=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- vSphere (cloud-config)"
bosh interpolate vsphere/cloud-config.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v network_name=test \
  -v vcenter_cluster=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- Azure"
bosh interpolate bosh.yml \
  -o azure/cpi.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v vnet_name=test \
  -v subnet_name=test \
  -v subscription_id=test \
  -v tenant_id=test \
  -v client_id=test \
  -v client_secret=test \
  -v resource_group_name=test \
  -v storage_account_name=test \
  -v default_security_group=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- VirtualBox with BOSH Lite"
bosh interpolate bosh.yml \
  -o virtualbox/cpi.yml \
  -o bosh-lite.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=vbox \
  -v internal_ip=192.168.56.6 \
  -v internal_gw=192.168.56.1 \
  -v internal_cidr=192.168.56.0/24 \
  -v network_name=vboxnet0 \
  -v host_ip=172.20.10.4 \
  -v host_username=test \
  -v host_private_key=test \
  --var-errs \
  --var-errs-unused \
  > /dev/null

echo "- VirtualBox (cloud-config)"
bosh interpolate virtualbox/cloud-config.yml \
  --var-errs \
  --var-errs-unused \
  > /dev/null
