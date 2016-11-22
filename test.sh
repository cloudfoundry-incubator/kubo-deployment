#!/bin/bash

set -eu

vars_store_prefix=/tmp/bosh-deployment-test

clean_tmp() {
  rm -f ${vars_store_prefix}.*
}

trap clean_tmp EXIT

echo "- AWS"
bosh interpolate bosh.yml --var-errs \
  -o use-aws.yml \
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
  > /dev/null

echo "- AWS with UAA"
bosh interpolate bosh.yml --var-errs \
  -o use-aws.yml \
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
  > /dev/null

echo "- GCP"
bosh interpolate bosh.yml --var-errs \
  -o use-gcp.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v service_account=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  > /dev/null

echo "- GCP with UAA"
bosh interpolate bosh.yml --var-errs \
  -o use-gcp.yml \
  -o uaa.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v service_account=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  > /dev/null

echo "- GCP with BOSH Lite"
bosh interpolate bosh.yml --var-errs \
  -o use-gcp.yml \
  -o enable-bosh-lite.yml \
  --vars-store $(mktemp ${vars_store_prefix}.XXXXXX) \
  -v director_name=test \
  -v internal_cidr=test \
  -v internal_gw=test \
  -v internal_ip=test \
  -v service_account=test \
  -v project_id=test \
  -v zone=test \
  -v network=test \
  -v subnetwork=test \
  > /dev/null
