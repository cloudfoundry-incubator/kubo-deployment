#!/bin/bash

set -eu

cd ..

tmp_file=/tmp/bosh-deployment-test
touch $tmp_file

clean_tmp() {
  rm -f $tmp_file
  rm -f ${tmp_file}.*
}

trap clean_tmp EXIT

# Only used for tests below. Ignore it.
function bosh() {
  shift 1
  command bosh int --var-errs --var-errs-unused ${@//--state=*/} > /dev/null
}

echo -e "\nCheck YAML syntax\n"
find .|grep yml|xargs -n1 bosh int

echo -e "\nUsed compiled releases\n"
grep -r -i s3.amazonaws.com/bosh-compiled-release-tarballs . | grep -v grep | grep -v ./.git

echo -e "\nUsed stemcells\n"
grep -r -i d/stemcells . | grep -v grep | grep -v ./.git

echo -e "\nExamples\n"

echo "- AWS"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v subnet_id=test

echo "- AWS with UAA"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v subnet_id=test

echo "- AWS with UAA + config-server"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  -o misc/config-server.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v subnet_id=test

echo "- AWS with UAA + CredHub + Turbulence"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  -o credhub.yml \
  -o turbulence.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v credhub_encryption_password=test

echo "- AWS with UAA for BOSH development"
bosh deploy bosh.yml \
  -o aws/cpi.yml \
  -o uaa.yml \
  -o misc/bosh-dev.yml \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_ip=test \
  -v access_key_id=test \
  -v secret_access_key=test \
  -v region=test \
  -v default_key_name=test \
  -v default_security_groups=[test]

echo "- AWS with external db and dns"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  -o misc/external-db.yml \
  -o misc/dns.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v internal_dns=[8.8.8.8] \
  -v access_key_id=test \
  -v secret_access_key=test \
  -v az=test \
  -v region=test \
  -v default_key_name=test \
  -v default_security_groups=[test] \
  -v private_key=test \
  -v subnet_id=test \
  -v external_db_host=test \
  -v external_db_port=test \
  -v external_db_user=test \
  -v external_db_password=test \
  -v external_db_adapter=test \
  -v external_db_name=test

echo "- AWS (cloud-config)"
bosh update-cloud-config aws/cloud-config.yml \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v az=test \
  -v subnet_id=test

echo "- GCP"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test

echo "- GCP with UAA"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  -o uaa.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test

echo "- GCP with UAA on external IP"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  -o uaa.yml \
  -o external-ip-not-recommended.yml \
  -o external-ip-not-recommended-uaa.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test \
  -v external_ip=test

echo "- GCP with BOSH Lite"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  -o bosh-lite.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test

echo "- GCP with BOSH Lite on Docker"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  -o bosh-lite-docker.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test

echo "- GCP with external db"
bosh create-env bosh.yml \
  -o gcp/cpi.yml \
  -o misc/external-db.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v gcp_credentials_json=test \
  -v project_id=test \
  -v zone=test \
  -v tags=[internal,no-ip] \
  -v network=test \
  -v subnetwork=test \
  -v external_db_host=test \
  -v external_db_port=test \
  -v external_db_user=test \
  -v external_db_password=test \
  -v external_db_adapter=test \
  -v external_db_name=test

echo "- GCP (cloud-config)"
bosh update-cloud-config gcp/cloud-config.yml \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  -v tags=[tag]

echo "- Openstack"
bosh create-env bosh.yml \
  -o openstack/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v openstack_domain=test \
  -v openstack_project=test \
  -v private_key=test \
  -v region=test

echo "- Openstack (cloud-config)"
bosh update-cloud-config openstack/cloud-config.yml \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v az=test \
  -v net_id=test

echo "- vSphere"
bosh create-env bosh.yml \
  -o vsphere/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v vcenter_cluster=test

echo "- vCloud"
bosh create-env bosh.yml \
  -o vcloud/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v network_name=test \
  -v vcloud_url=test \
  -v vcloud_user=test \
  -v vcloud_password=test \
  -v vcd_org=test \
  -v vcd_name=test

echo "- vSphere (cloud-config)"
bosh update-cloud-config vsphere/cloud-config.yml \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v network_name=test \
  -v vcenter_cluster=test

echo "- Azure"
bosh create-env bosh.yml \
  -o azure/cpi.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=10.0.0.0/24 \
  -v internal_gw=10.0.0.1 \
  -v internal_ip=10.0.0.4 \
  -v vnet_name=boshvnet-crp \
  -v subnet_name=Bosh \
  -v subscription_id=test \
  -v tenant_id=test \
  -v client_id=test \
  -v client_secret=test \
  -v resource_group_name=test \
  -v storage_account_name=test \
  -v default_security_group=nsg-bosh

echo "- Azure (custom-environment)"
bosh create-env bosh.yml \
  -o azure/cpi.yml \
  -o azure/custom-environment.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=10.0.0.0/24 \
  -v internal_gw=10.0.0.1 \
  -v internal_ip=10.0.0.4 \
  -v vnet_name=boshvnet-crp \
  -v subnet_name=Bosh \
  -v environment=AzureChinaCloud \
  -v subscription_id=test \
  -v tenant_id=test \
  -v client_id=test \
  -v client_secret=test \
  -v resource_group_name=test \
  -v storage_account_name=test \
  -v default_security_group=nsg-bosh

echo "- Azure (cloud-config)"
bosh update-cloud-config azure/cloud-config.yml \
  -v internal_cidr=10.0.16.0/24 \
  -v internal_gw=10.0.16.1 \
  -v vnet_name=boshvnet-crp \
  -v subnet_name=CloudFoundry \
  -v security_group=nsg-cf

echo "- VirtualBox with BOSH Lite"
bosh create-env bosh.yml \
  -o virtualbox/cpi.yml \
  -o bosh-lite.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=vbox \
  -v internal_ip=192.168.56.6 \
  -v internal_gw=192.168.56.1 \
  -v internal_cidr=192.168.56.0/24

echo "- VirtualBox with IPv6 (remote)"
bosh create-env bosh.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -o virtualbox/cpi.yml \
  -o virtualbox/outbound-network.yml \
  -o jumpbox-user.yml \
  -o uaa.yml \
  -o credhub.yml \
  -o misc/ipv6/bosh.yml \
  -o misc/ipv6/uaa.yml \
  -o misc/ipv6/credhub.yml \
  -o virtualbox/remote.yml \
  -o virtualbox/ipv6/cpi.yml \
  -o virtualbox/ipv6/remote.yml \
  -v director_name=vbox \
  -v internal_cidr=fd7a:eeed:e696:969f:0000:0000:0000:0000/64 \
  -v internal_gw=fd7a:eeed:e696:969f:0000:0000:0000:0001 \
  -v internal_ip=fd7a:eeed:e696:969f:0000:0000:0000:0004 \
  -v outbound_network_name=NatNetwork \
  -v vbox_host=fd7a:eeed:e696:969f:0000:0000:0000:0001 \
  -v vbox_username=test

echo "- VirtualBox with BOSH Lite with garden-runc"
bosh create-env bosh.yml \
  -o virtualbox/cpi.yml \
  -o bosh-lite.yml \
  -o bosh-lite-runc.yml \
  -o jumpbox-user.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=vbox \
  -v internal_ip=192.168.56.6 \
  -v internal_gw=192.168.56.1 \
  -v internal_cidr=192.168.56.0/24

echo "- Warden (cloud-config)"
bosh update-cloud-config warden/cloud-config.yml

echo "- Docker"
bosh create-env bosh.yml \
  -o docker/cpi.yml \
  -o jumpbox-user.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=docker \
  -v internal_cidr=10.245.0.0/16 \
  -v internal_gw=10.245.0.1 \
  -v internal_ip=10.245.0.10 \
  -v docker_host=tcp://192.168.50.8:4243 \
  --var-file docker_tls.ca=$tmp_file \
  --var-file docker_tls.certificate=$tmp_file \
  --var-file docker_tls.private_key=$tmp_file \
  -v network=net3

echo "- Docker via UNIX sock"
bosh create-env bosh.yml \
  -o docker/cpi.yml \
  -o docker/unix-sock.yml \
  -o jumpbox-user.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=docker \
  -v internal_cidr=10.245.0.0/16 \
  -v internal_gw=10.245.0.1 \
  -v internal_ip=10.245.0.10 \
  -v docker_host=unix:///var/run/docker.sock \
  -v network=net3

echo "- Docker (cloud-config)"
bosh update-cloud-config docker/cloud-config.yml -v network=net3

echo "- Warden"
bosh create-env bosh.yml \
  -o warden/cpi.yml \
  -o jumpbox-user.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
  -v director_name=docker \
  -v internal_cidr=10.245.0.0/16 \
  -v internal_gw=10.245.0.1 \
  -v internal_ip=10.245.0.10 \
  -v garden_host=127.0.0.1

echo "- Secondary CPIs"
bosh create-env bosh.yml \
  -o aws/cpi.yml \
  -o docker/cpi-secondary.yml \
  -o azure/cpi-secondary.yml \
  -o vsphere/cpi-secondary.yml \
  -o openstack/cpi-secondary.yml \
  --state=$tmp_file \
  --vars-store $(mktemp ${tmp_file}.XXXXXX) \
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
  -v subnet_id=test
