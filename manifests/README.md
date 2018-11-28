# Deploying CFCR

The base manifest "just works" and will deploy a running cluster of Kubernetes:

```
bosh -d cfcr deploy kubo-deployment/manifests/cfcr.yml
```

For deeper documentation to deploy CFCR go [here](https://github.com/cloudfoundry-incubator/kubo-release/#deploying-cfcr).


## Operator files

### BOSH options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/use-runtime-config-bosh-dns.yml`](ops-files/use-runtime-config-bosh-dns.yml) |  | Deprecated, currently empty. To be removed in next release |
| [`ops-files/rename.yml`](ops-files/rename.yml) | Specify the deployment name | The deployment name is also used for etcd certificates. |
| [`ops-files/vm-types.yml`](ops-files/vm-types.yml) | Specify the `vm_type` for `master`, `worker` and `apply-addons` instances | By default, `master`, `worker` and `apply-addons` instances assume `vm_type: small`, `vm_type: small-highmem` and `vm_type: minimal`, respectively (`vm_types` that are also assumed to exists by https://github.com/cloudfoundry/cf-deployment manifests). You may want to use bespoke `vm_types` so as to scale them, tag them, or apply unique `cloud_properties` independently of other deployments in the same BOSH environment. |
| [`ops-files/add-vm-extensions-to-master.yml`](ops-files/add-vm-extensions-to-master.yml) | Add VM Extensions for loadbalancers to master | |
| [`ops-files/use-vm-extensions.yml`](ops-files/use-vm-extensions.yml) | Configure the `master` and `worker` instance groups on AWS and GCP to consume their respective `vm_extensions` | Only works when used in tandem with the BOSH cloud-configs for AWS or GCP outlined below |
| [`ops-files/iaas/vsphere/use-vm-extensions.yml`](ops-files/iaas/vsphere/use-vm-extensions.yml) | Configure vSphere `worker` instance groups to consume their respective `vm_extensions` | Only works when used in tandem with the BOSH cloud-config for vSphere outlined below |
| [`ops-files/worker_count.yml`](ops-files/worker_count.yml) | Specify the count for `worker` instances | By default, 3 `worker` instances. |
| [`ops-files/non-precompiled-releases.yml`](ops-files/non-precompiled-releases.yml) | Use non-precompiled releases when deploying CFCR. |

### BOSH Cloud Config

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| **AWS**
| [`cloud-config/iaas/aws/use-vm-extensions.yml`](cloud-config/iaas/aws/use-vm-extensions.yml) | Configure the cloud-config to control the AWS Cloud Provider using `vm_extensions`  | |
| **GCP**
| [`cloud-config/iaas/gcp/use-vm-extensions.yml`](cloud-config/iaas/gcp/use-vm-extensions.yml) | Configure the cloud-config to control the GCP Cloud Provider using `vm_extensions` | |
| **vSphere**
| [`cloud-config/iaas/vsphere/use-vm-extensions.yml`](cloud-config/iaas/vsphere/use-vm-extensions.yml) | Configure the cloud-config to control the vSphere Cloud Provider using `vm_extensions` | |

### Routing options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| **OpenStack** | | |
| [`ops-files/iaas/openstack/master-static-ip.yml`](ops-files/iaas/openstack/master-static-ip.yml) | Attach floating IP to Kube API | Assign allocated floating IP to `master` instance. IP included in TLS certificates. |
| **vSphere** | | |
| [`ops-files/iaas/vsphere/master-static-ip.yml`](ops-files/iaas/vsphere/master-static-ip.yml) | Assign static IP to Kube API | Assign static IP to `master` instance. IP included in TLS certificates. |
| **gcp** | | |
| [`ops-files/iaas/gcp/add-service-key-master.yml`](ops-files/iaas/gcp/add-service-key-master.yml) | Allow user to specify GCP key instead of service account |  |
| [`ops-files/iaas/gcp/add-service-key-worker.yml`](ops-files/iaas/gcp/add-service-key-worker.yml) | Allow user to specify GCP key instead of service account |  |

### Infrastructure

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| **AWS** | | |
| [`ops-files/iaas/aws/cloud-provider.yml`](ops-files/iaas/aws/cloud-provider.yml) | Enable Cloud Provider for AWS | Requires AWS Instance Profiles (not API keys) to grant Kubernetes access to AWS |
| [`ops-files/iaas/aws/lb.yml`](ops-files/iaas/aws/lb.yml) | Enable instance tagging for AWS |  |
| [`ops-files/iaas/aws/add-master-credentials.yml`](ops-files/iaas/aws/add-master-credentials.yml) | Set AWS credentials for the Kube API and Kube Controller Manager |  |
| [`ops-files/iaas/aws/add-worker-credentials.yml`](ops-files/iaas/aws/add-worker-credentials.yml) | Set AWS credentials for the Kubelet |  |
| **OpenStack** | | |
| [`ops-files/iaas/openstack/cloud-provider.yml`](ops-files/iaas/openstack/cloud-provider.yml) | Enable Cloud Provider for OpenStack | Enable Cloud Provider for OpenStack |
| **GCP** | | |
| [`ops-files/iaas/gcp/cloud-provider.yml`](ops-files/iaas/gcp/cloud-provider.yml) | Enable Cloud Provider for GCP | - |
| [`ops-files/iaas/gcp/add-subnetwork-for-internal-load-balancer.yml`](ops-files/iaas/gcp/add-subnetwork-for-internal-load-balancer.yml) | Specify subnetwork for GCP | Cloud Provider has to be enabled first |
| **vSphere** | | |
| [`ops-files/iaas/vsphere/cloud-provider.yml`](ops-files/iaas/vsphere/cloud-provider.yml) | Enable Cloud Provider for vSphere | - |
| [`ops-files/iaas/vsphere/set-working-dir-no-rp.yml`](ops-files/iaas/vsphere/set-working-dir-no-rp.yml) | Configure vSphere cloud provider's working dir if there is no resource pool | - |
| **virtualbox** | | |
| [`ops-files/iaas/virtualbox/bosh-lite.yml`](ops-files/iaas/virtualbox/bosh-lite.yml) | Enables CFCR to run on a virtualbox bosh-lite environment | Deploys 1 master and 3 workers. Master is deployed to a static ip: 10.244.0.34 |
| **Azure** | | |
| [`ops-files/iaas/azure/cloud-provider.yml`](ops-files/iaas/azure/cloud-provider.yml) | Enable Cloud Provider for Azure | Requires Azure CPI >= v35.5.0 |
| [`ops-files/iaas/azure/subnet.yml`](ops-files/iaas/azure/subnet.yml) | Changes the subnet | |
| [`ops-files/iaas/azure/use-cifs.yml`](ops-files/iaas/azure/use-cifs.yml) | Installs CIFS utils and allows using azure-file volume | |
| [`ops-files/iaas/azure/use-credentials.yml`](ops-files/iaas/azure/use-credentials.yml) | Uses AD credentials instead of Managed Identity | |

### Proxy

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/add-proxy.yml`](ops-files/add-proxy.yml) | Configure HTTP_PROXY, HTTPS_PROXY, and NO_PROXY for Kubernetes components | All Kubernetes components are configured with the `http_proxy`, `https_proxy`, and `no_proxy` environment variables |

### Kubernetes

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/addons-spec.yml`](ops-files/addons-spec.yml) | Addons to be deployed into the Kubernetes cluster | - |
| [`ops-files/allow-privileged-containers.yml`](ops-files/allow-privileged-containers.yml) | Allows privileged containers for the Kubernetes cluster. | It is not recommended to use privileged containers however some workloads require it. Container privileges can be limited with the SecurityContextDeny admission plugin (set by default in CFCR). See kubernetes documentation for more information |
| [`ops-files/disable-anonymous-auth.yml`](ops-files/disable-anonymous-auth.yml) | Disable `anonymous-auth` on the API server | - |
| [`ops-files/add-oidc-endpoint.yml`](ops-files/add-oidc-endpoint.yml) | Enable OIDC authentication for the Kubernetes cluster | - |
| [`ops-files/change-cidrs.yml`](ops-files/change-cidrs.yml) | Change POD CIDR and Service Cluster CIDR. This should only be applied to a new cluster, please do not apply to an existing cluster. | Extra Vars Required:<br>- **first_ip_of_service_cluster_cidr:** Required for TLS certificate of apiserver<br>- **kubedns_service_ip**: Required for kube dns IP address, needs to be part of service_cluster_cidr |
| [`ops-files/enable-denyescalatingexec.yml`](ops-files/enable-denyescalatingexec.yml) | Enables the DenyEscalatingExec admission plugin. | This ops-file is recommended for most clusters. | - |
| [`ops-files/enable-securitycontextdeny.yml`](ops-files/enable-securitycontextdeny.yml) | Enables the SecurityContextDeny admission plugin. | This ops-file is recommended for most clusters. | - |
| [`ops-files/enable-podsecuritypolicy.yml`](ops-files/enable-podsecuritypolicy.yml) | Enables the PodSecurityPolicy admission plugin. | Please ensure that you have applied an appropriate policy before enabling this plugin.  Failure to do so will result in failure of your workloads. |
| [`ops-files/add-hostname-to-master-certificate.yml`](ops-files/add-hostname-to-master-certificate.yml) | Add hostname to master certificate | Extra Vars Required:<br>- **api-hostname:** Required for TLS certificate of apiserver |
| [`ops-files/use-coredns.yml`](ops-files/use-coredns.yml) | Add CoreDNS to the list of addons deployed by the apply-specs errand | - |
| [`ops-files/enable-encryption-config.yml`](ops-files/enable-encryption-config.yml) | Enable data encryption at rest | Extra Vars Required:<br>- **encryption-config:** Encryption configuration as described [here](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/#understanding-the-encryption-at-rest-configuration). Must be a file, interpolated with `--var-file`. Example: `--var-file encryption-config=encryption-config.yml`  |

### Etcd

| Name | Purpose | Notes|
|:--- |:--- |:--- |
| [`ops-files/change-etcd-metrics-url.yml`](ops-files/change-etcd-metrics-url.yml) | Change procotol and port of the etcd's metrics endpoint | - |

### Certificates

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/set-certificate-duration.yml`](ops-files/set-certificate-duration.yml) | Set the duration of all generated certificates to a specified duration | Extra Vars Required:<br>- certificate-duration: Duration, specified in days, for all certificates generated in manifest |

### BOSH Backup & Restore

| Name | Purpose | Notes|
|:--- |:--- |:--- |
| [`ops-files/enable-bbr.yml`](ops-files/enable-bbr.yml) | Deploy jobs required to enable BBR. | Only tested with single master. |

### NFS

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/enable-nfs.yml`](ops-files/enable-nfs.yml) | Enables packages to be install on worker vms required for NFS | - |

### Dev

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/kubo-local-release.yml`](ops-files/kubo-local-release.yml) | Deploy a local kubo release located in `../kubo-release` | -  |
