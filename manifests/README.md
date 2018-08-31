# Deploying CFCR

The base manifest "just works" and will deploy a running cluster of Kubernetes:

```
bosh -d cfcr deploy kubo-deployment/manifests/cfcr.yml
```

## Dependencies

The only dependencies are that your BOSH environment has:

* Credhub/UAA
* Cloud Config with `vm_types` named `minimal`, `small`, and `small-highmem` as per similar requirements of [cf-deployment](https://github.com/cloudfoundry/cf-deployment)
* Cloud Config has a network named `default`as per similar requirements of [cf-deployment](https://github.com/cloudfoundry/cf-deployment)
* If using bosh-lite see [Deploy CFCR in bosh-lite](https://github.com/cloudfoundry-incubator/kubo-deployment/blob/master/CONTRIBUTING.md#deploy-cfcr-in-bosh-lite)
* Ubuntu Trusty stemcell `3468` is uploaded (it's up to you to keep up to date with latest `3468.X` versions and update your BOSH deployments)

## Getting Started

You can get started with one `bosh deploy` command. It will download and deploy everything for you.

```
export BOSH_ENVIRONMENT=<bosh-name>
export BOSH_DEPLOYMENT=cfcr
git clone https://github.com/cloudfoundry-incubator/kubo-deployment
bosh deploy kubo-deployment/manifests/cfcr.yml
```

To see the running cluster:

```
$ bosh instances

Deployment 'cfcr'

Instance                                     Process State  AZ  IPs
master/bde7bc5a-a9fd-4bcc-9ba7-b66752fad159  running        z1  10.10.1.20
worker/4518c694-3328-4538-bc08-dedf8a3bf400  running        z1  10.10.1.22
worker/49d317d0-dff2-44a3-b00c-0406ce8a010e  running        z1  10.10.1.23
worker/e00ac851-fadb-4b7d-94c4-8917042ba6cb  running        z1  10.10.1.21
```

Once the deployment is running, you can setup your `kubectl` CLI to connect and authenticate you.

First, get the randomly generated Kubernetes API admin password from Credhub:

```
admin_password=$(bosh int <(credhub get -n "${BOSH_ENVIRONMENT}/${BOSH_DEPLOYMENT}/kubo-admin-password" --output-json) --path=/value)
```

Next, get the dynamically assigned IP address of the `master/0` instance:

```
master_host=$(bosh int <(bosh instances --json) --path /Tables/0/Rows/0/ips)
```

Finally, setup your local `kubectl` configuration:

```
cluster_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}"
user_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}-admin"
context_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}"

kubectl config set-cluster "${cluster_name}" \
  --server="https://${master_host}:8443" \
  --insecure-skip-tls-verify=true
kubectl config set-credentials "${user_name}" --token="${admin_password}"
kubectl config set-context "${context_name}" --cluster="${cluster_name}" --user="${user_name}"
kubectl config use-context "${context_name}"
```

To confirm that you are connected and configured to your Kubernetes cluster:

```
$ kubectl get all
NAME             TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
svc/kubernetes   ClusterIP   10.100.200.1   <none>        443/TCP   2h
```


## Operator files

### BOSH options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/use-trusty-stemcell.yml`](ops-files/use-trusty-stemcell.yml) | (Deprecated) Use trusty stemcell |  |
| [`ops-files/use-runtime-config-bosh-dns.yml`](ops-files/use-runtime-config-bosh-dns.yml) | Delegate `bosh-dns` addon to BOSH runtime config | Apply this operator file if your BOSH environment has a runtime config that adds the `bosh-dns` job to all instances. By default, `cfcr.yml` will add `bosh-dns` to deployment instances. |
| [`ops-files/rename.yml`](ops-files/rename.yml) | Specify the deployment name | The deployment name is also used for etcd certificates. |
| [`ops-files/vm-types.yml`](ops-files/vm-types.yml) | Specify the `vm_type` for `master`, `worker` and `apply-addons` instances | By default, `master`, `worker` and `apply-addons` instances assume `vm_type: small`, `vm_type: small-highmem` and `vm_type: minimal`, respectively (`vm_types` that are also assumed to exists by https://github.com/cloudfoundry/cf-deployment manifests). You may want to use bespoke `vm_types` so as to scale them, tag them, or apply unique `cloud_properties` independently of other deployments in the same BOSH environment. |
| [`ops-files/add-vm-extensions-to-master.yml`](ops-files/add-vm-extensions-to-master.yml) | Add VM Extensions for loadbalancers to master | |
| [`ops-files/use-vm-extensions.yml`](ops-files/use-vm-extensions.yml) | Configure the `master` and `worker` instance groups to consume their respective `vm_extensions` | Only works when used in tandem with the BOSH cloud-configs outlined below |

### BOSH Cloud Config

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| **AWS**
| [`cloud-config/iaas/aws/use-vm-extensions.yml`](cloud-config/iaas/aws/use-vm-extensions.yml) | Configure the cloud-config to control the AWS Cloud Provider using `vm_extensions`  | |
| **GCP**
| [`cloud-config/iaas/gcp/use-vm-extensions.yml`](cloud-config/iaas/gcp/use-vm-extensions.yml) | Configure the cloud-config to control the GCP Cloud Provider using `vm_extensions` | |

### Routing options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/cf-routing.yml`](ops-files/cf-routing.yml) | Combine CFCR with Cloud Foundry routing | (Deprecated) Kube API and labeled services advertised to Cloud Foundry routing mesh. Kube API hostname is included in TLS certificates. |
| [`ops-files/cf-routing-links.yml`](ops-files/cf-routing-links.yml) | As above, using BOSH links | (Deprecated) Simpler method of integration with Cloud Foundry running on same BOSH [in development] |
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

### Proxy

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/add-proxy.yml`](ops-files/add-proxy.yml) | Configure HTTP_PROXY, HTTPS_PROXY, and NO_PROXY for Kubernetes components | All Kubernetes components are configured with the `http_proxy`, `https_proxy`, and `no_proxy` environment variables |

### Kubernetes

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/addons-spec.yml`](ops-files/addons-spec.yml) | Addons to be deployed into the Kubernetes cluster | - |
| [`ops-files/allow-privileged-containers.yml`](ops-files/allow-privileged-containers.yml) | Allows privileged containers for the Kubernetes cluster | - |
| [`ops-files/disable-anonymous-auth.yml`](ops-files/disable-anonymous-auth.yml) | Disable `anonymous-auth` on the API server | - |
| [`ops-files/disable-deny-escalating-exec.yml`](ops-files/disable-deny-escalating-exec.yml) | Disable `DenyEscalatingExec` in API server admission control | - |
| [`ops-files/add-oidc-endpoint.yml`](ops-files/add-oidc-endpoint.yml) | Enable OIDC authentication for the Kubernetes cluster | - |
| [`ops-files/change-cidrs.yml`](ops-files/change-cidrs.yml) | Change POD CIDR and Service Cluster CIDR. This should only be applied to a new cluster, please do not apply to an existing cluster. | Extra Vars Required:<br>- **first_ip_of_service_cluster_cidr:** Required for TLS certificate of apiserver<br>- **kubedns_service_ip**: Required for kube dns IP address, needs to be part of service_cluster_cidr |
| [`ops-files/add-hostname-to-master-certificate.yml`](ops-files/add-hostname-to-master-certificate.yml) | Add hostname to master certificate | Extra Vars Required:<br>- **api-hostname:** Required for TLS certificate of apiserver |
| [`ops-files/use-coredns.yml`](ops-files/use-coredns.yml) | Add CoreDNS to the list of addons deployed by the apply-specs errand | - |

### BOSH Backup & Restore (Experimental)

| Name | Purpose | Notes|
|:--- |:--- |:--- |
| [`ops-files/enable-bbr.yml`](ops-files/enable-bbr.yml) | Deploy jobs required to enable BBR. | Only tested with single master. |

### Dev

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`ops-files/kubo-local-release.yml`](ops-files/kubo-local-release.yml) | Deploy a local kubo release located in `../kubo-release` | -  |


