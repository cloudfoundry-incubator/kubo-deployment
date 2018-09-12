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
      check_interpolation "rename.yml" "-v deployment_name=fubar"
      check_interpolation "vm-types.yml" "-v master_vm_type=master" "-v worker_vm_type=worker" "-v apply_addons_vm_type=addons"
      check_interpolation "use-trusty-stemcell.yml"
      check_interpolation "add-vm-extensions-to-master.yml"
      check_interpolation "use-vm-extensions.yml" "-v deployment_name=cfcr"

      # Infrastructure
      check_interpolation "iaas/aws/cloud-provider.yml"
      check_interpolation "iaas/aws/lb.yml" "-v kubernetes_cluster_tag=test"
      check_interpolation "name:iaas/aws/add-master-credentials.yml" "iaas/aws/cloud-provider.yml" "-o iaas/aws/add-master-credentials.yml" "-v aws_access_key_id_master=access-key-id" "-v aws_secret_access_key_master=secret-access-key"
      check_interpolation "name:iaas/aws/add-worker-credentials.yml" "iaas/aws/cloud-provider.yml" "-o iaas/aws/add-worker-credentials.yml" "-v aws_access_key_id_worker=access-key-id" "-v aws_secret_access_key_worker=secret-access-key"
      check_interpolation "iaas/gcp/cloud-provider.yml" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "name:iaas/gcp/add-subnetwork-for-internal-load-balancer.yml" "iaas/gcp/cloud-provider.yml" "-o iaas/gcp/add-subnetwork-for-internal-load-balancer.yml" "-v subnetwork=foo" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "name:iaas/gcp/add-service-key-master.yml" "iaas/gcp/cloud-provider.yml" "-o iaas/gcp/add-service-key-master.yml" "-v service_key_master=foo" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "name:iaas/gcp/add-service-key-worker.yml" "iaas/gcp/cloud-provider.yml" "-o iaas/gcp/add-service-key-worker.yml" "-v service_key_worker=foo" "-l example-vars-files/iaas/gcp/cloud-provider.yml"
      check_interpolation "iaas/openstack/master-static-ip.yml" "-v kubernetes_master_host=10.11.12.13"
      check_interpolation "iaas/openstack/cloud-provider.yml" "-l example-vars-files/iaas/openstack/cloud-provider.yml"
      check_interpolation "iaas/vsphere/cloud-provider.yml" "-l example-vars-files/iaas/vsphere/cloud-provider.yml"
      check_interpolation "name:iaas/vsphere/set-working-dir-no-rp.yml" "iaas/vsphere/cloud-provider.yml" "-o iaas/vsphere/set-working-dir-no-rp.yml" "-l example-vars-files/iaas/vsphere/set-working-dir-no-rp.yml"
      check_interpolation "iaas/vsphere/master-static-ip.yml" "-v kubernetes_master_host=10.11.12.13"
      check_interpolation "iaas/vsphere/use-vm-extensions.yml"
      check_interpolation "iaas/virtualbox/bosh-lite.yml"

      # Routing Variations
      check_interpolation "cf-routing.yml" "-l example-vars-files/cf-routing.yml"
      check_interpolation "cf-routing-links.yml" "-l example-vars-files/cf-routing-links.yml"

      # HTTP proxy options
      check_interpolation "add-proxy.yml" "-v http_proxy=10.10.10.10:8000 -v https_proxy=10.10.10.10:8000 -v no_proxy=localhost,127.0.0.1"

      # Syslog
      check_interpolation "add-syslog.yml" "-l example-vars-files/add-syslog.yml"  
      check_interpolation "name:add-syslog-tls.yml" "add-syslog.yml" "-o add-syslog-tls.yml" "-l example-vars-files/add-syslog.yml" "-l example-vars-files/add-syslog-tls.yml"

      # Kubernetes
      check_interpolation "addons-spec.yml" "-v addons-spec={}"
      check_interpolation "allow-privileged-containers.yml"
      check_interpolation "disable-anonymous-auth.yml"
      check_interpolation "disable-deny-escalating-exec.yml"
      check_interpolation "add-oidc-endpoint.yml" "-l example-vars-files/misc/oidc.yml"
      check_interpolation "change-cidrs.yml" "-l example-vars-files/new-cidrs.yml"
      check_interpolation "add-hostname-to-master-certificate.yml" "-v api-hostname=example.com"
      check_interpolation "use-coredns.yml"

      # BBR
      check_interpolation "enable-bbr.yml"

      # Dev
      check_interpolation "kubo-local-release.yml"

    popd > /dev/null # operations
  popd > /dev/null
  exit $exit_code
}
