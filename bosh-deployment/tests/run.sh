#!/bin/bash

set -e # -x

tests_dir=$PWD

cd ..

rm -f $tests_dir/creds.yml

echo "-----> `date`: Create env"
bosh create-env bosh.yml \
  --state $tests_dir/state.json \
  -o virtualbox/cpi.yml \
  -o virtualbox/outbound-network.yml \
  -o bosh-lite.yml \
  -o bosh-lite-runc.yml \
  -o jumpbox-user.yml \
  --vars-store $tests_dir/creds.yml \
  -v director_name="Bosh Lite Director" \
  -v internal_ip=192.168.50.10 \
  -v internal_gw=192.168.50.1 \
  -v internal_cidr=192.168.50.0/24 \
  -v outbound_network_name=NatNetwork

export BOSH_ENVIRONMENT=192.168.50.10
export BOSH_CA_CERT="$(bosh int $tests_dir/creds.yml --path /director_ssl/ca)"
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET="$(bosh int $tests_dir/creds.yml --path /admin_password)"

echo "-----> `date`: Update cloud config"
bosh -n update-cloud-config warden/cloud-config.yml

echo "-----> `date`: Upload stemcell"
bosh -n upload-stemcell "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3468.1" \
  --sha1 69bbf7a8c4683a8130eaf22b3270f8d737037884

echo "-----> `date`: Deploy"
bosh -n -d zookeeper deploy <(wget -O- https://raw.githubusercontent.com/cppforlife/zookeeper-release/master/manifests/zookeeper.yml)

echo "-----> `date`: Exercise deployment"
bosh -n -d zookeeper run-errand smoke-tests

echo "-----> `date`: Clean up disks, etc."
bosh -n -d zookeeper clean-up --all

echo "-----> `date`: Deleting env"
bosh delete-env bosh.yml \
  --state $tests_dir/state.json \
  -o virtualbox/cpi.yml \
  -o virtualbox/outbound-network.yml \
  -o bosh-lite.yml \
  -o bosh-lite-runc.yml \
  -o jumpbox-user.yml \
  --vars-store $tests_dir/creds.yml \
  -v director_name="Bosh Lite Director" \
  -v internal_ip=192.168.50.10 \
  -v internal_gw=192.168.50.1 \
  -v internal_cidr=192.168.50.0/24 \
  -v outbound_network_name=NatNetwork

echo "-----> `date`: Done"
