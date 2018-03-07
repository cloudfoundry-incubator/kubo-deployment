#!/bin/bash

test_standard_ops() {
  # Padded for pretty output
  suite_name="STANDARD    "

  pushd ${home}/manifests > /dev/null
    pushd ops-files > /dev/null
      if interpolate ""; then
        pass "cfcr.yml"
      else
        fail "cfcr.yml"
      fi

      # CI & wrapper scripts
      check_interpolation "misc/bootstrap.yml" "-l example-vars-files/misc/bootstrap.yml"
      check_interpolation "misc/bootstrap.yml"  "-o misc/dev.yml"  "-o use-runtime-config-bosh-dns.yml" "-l example-vars-files/misc/bootstrap.yml"

      # BOSH
      check_interpolation "use-runtime-config-bosh-dns.yml"
      check_interpolation "vm-types.yml" "-v master_vm_type=master" "-v worker_vm_type=worker"

      # Infrastructure
      check_interpolation "iaas/aws/cloud-provider.yml"
      check_interpolation "iaas/aws/lb.yml" "-v kubernetes_cluster_tag=test"
      check_interpolation "name:iaas/aws/add-master-credentials.yml" "iaas/aws/cloud-provider.yml" "-o iaas/aws/add-master-credentials.yml" "-v aws_access_key_id_master=access-key-id" "-v aws_secret_access_key_master=secret-access-key"
      check_interpolation "name:iaas/aws/add-worker-credentials.yml" "iaas/aws/cloud-provider.yml" "-o iaas/aws/add-worker-credentials.yml" "-v aws_access_key_id_worker=access-key-id" "-v aws_secret_access_key_worker=secret-access-key"
      check_interpolation "iaas/gcp/cloud-provider.yml" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "name:iaas/gcp/add-service-key-master.yml" "iaas/gcp/cloud-provider.yml" "-o iaas/gcp/add-service-key-master.yml" "-v service_key_master=foo" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "name:iaas/gcp/add-service-key-worker.yml" "iaas/gcp/cloud-provider.yml" "-o iaas/gcp/add-service-key-worker.yml" "-v service_key_worker=foo" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "iaas/openstack/master-static-ip.yml" "-v kubernetes_master_host=10.11.12.13"
      check_interpolation "iaas/vsphere/cloud-provider.yml" "-l example-vars-files/iaas/vsphere/cloud-provider.yml"
      check_interpolation "iaas/vsphere/cloud-provider.yml" "-o iaas/vsphere/set-working-dir-no-rp.yml" "-l example-vars-files/iaas/vsphere/set-working-dir-no-rp.yml"
      check_interpolation "iaas/vsphere/master-static-ip.yml" "-v kubernetes_master_host=10.11.12.13"

      # Routing Variations
      check_interpolation "cf-routing.yml" "-l example-vars-files/cf-routing.yml"
      check_interpolation "cf-routing-links.yml" "-l example-vars-files/cf-routing-links.yml"

      # HTTP proxy options
      check_interpolation "add-proxy.yml" "-v http_proxy=10.10.10.10:8000 -v https_proxy=10.10.10.10:8000 -v no_proxy=localhost,127.0.0.1"

      # Kubernetes
      check_interpolation "system-specs.yml" "-l example-vars-files/system-specs.yml"
      check_interpolation "name:addons-spec.yml" "system-specs.yml" "-o addons-spec.yml" "-v addons-spec={}" "-l example-vars-files/system-specs.yml"
      check_interpolation "allow-privileged-containers.yml"
      check_interpolation "add-oidc-endpoint.yml" "-l example-vars-files/misc/oidc.yml"

    popd > /dev/null # operations
  popd > /dev/null
  exit $exit_code
}
