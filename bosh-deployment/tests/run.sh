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
  -o uaa.yml \
  -o credhub.yml \
  -o jumpbox-user.yml \
  -o misc/blobstore-tls.yml \
  -o misc/nats-strict-tls.yml \
  --vars-store $tests_dir/creds.yml \
  -v director_name=bosh-lite \
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
bosh -n upload-stemcell "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3586.16" \
  --sha1 7b4b314abf3a8f06973f3533162be13d57ebed28

echo "-----> `date`: Deploy"
bosh -n -d zookeeper deploy <(wget -O- https://raw.githubusercontent.com/cppforlife/zookeeper-release/master/manifests/zookeeper.yml) \
  -o tests/cred-test.yml

echo "-----> `date`: Exercise deployment"
bosh -n -d zookeeper run-errand smoke-tests

echo "-----> `date`: Exercise deployment"
bosh -n -d zookeeper recreate

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
  -v director_name=bosh-lite \
  -v internal_ip=192.168.50.10 \
  -v internal_gw=192.168.50.1 \
  -v internal_cidr=192.168.50.0/24 \
  -v outbound_network_name=NatNetwork

echo "-----> `date`: Done"
