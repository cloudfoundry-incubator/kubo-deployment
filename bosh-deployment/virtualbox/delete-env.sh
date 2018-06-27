#!/bin/bash

set -eu

STEP() { echo ; echo ; echo "==\\" ; echo "===>" "$@" ; echo "==/" ; echo ; }

echo "This will destroy BOSH from VirtualBox."
echo

read -p "Continue? [yN] "
[[ $REPLY =~ ^[Yy]$ ]] || exit 1


####
STEP "Deleting BOSH Director"
####

bosh delete-env bosh-deployment/bosh.yml \
  --state state.json \
  --ops-file bosh-deployment/virtualbox/cpi.yml \
  --ops-file bosh-deployment/virtualbox/outbound-network.yml \
  --ops-file bosh-deployment/bosh-lite.yml \
  --ops-file bosh-deployment/bosh-lite-runc.yml \
  --ops-file bosh-deployment/jumpbox-user.yml \
  --vars-store creds.yml \
  --var director_name=bosh-lite \
  --var internal_ip=192.168.50.6 \
  --var internal_gw=192.168.50.1 \
  --var internal_cidr=192.168.50.0/24 \
  --var outbound_network_name=NatNetwork
