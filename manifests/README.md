# Deploying CFCR

The base manifest "just works" and will deploy a running cluster of Kubernetes:

```
bosh deploy kubo-deployment/manifests/cfcr.yml
```

## Dependencies

The only dependencies are that your BOSH environment has:

* Credhub/UAA
* Cloud Config with `vm_types` named `minimal`, `small`, and `small-highmem` as per similar requirements of [cf-deployment](https://github.com/cloudfoundry/cf-deployment)
* Cloud Config has a network named `default`as per similar requirements of [cf-deployment](https://github.com/cloudfoundry/cf-deployment)
* Not a bosh-lite
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

## Integrate with Cloud Foundry TCP routing

If you are already running Cloud Foundry, then you can reuse its TCP router to provide public access to your Kubernetes services and the Kubernetes API.

We will also be changing how we interact with the Kubernetes API. Instead of using https://IP:8443 we will access it through the Cloud Foundry TCP routing hostname and the selected port; such as https://tcp.mycompany.com:8443

So we need need to have some certificates regenerated to include the new hostname. Delete them from Credhub:

```
export BOSH_ENVIRONMENT=<bosh-name>
export BOSH_DEPLOYMENT=cfcr
credhub delete -n /$BOSH_ENVIRONMENT/$BOSH_DEPLOYMENT/tls-kubernetes
credhub delete -n /$BOSH_ENVIRONMENT/$BOSH_DEPLOYMENT/tls-kubelet
```

Next, we need to document information about your Cloud Foundry and how CFCR will be allowed to register TCP routes.

Create a file `cf-vars.yml` which might look like:

```yaml
kubernetes_master_host: tcp.apps.mycompany.com
kubernetes_master_port: 8443
routing_cf_api_url: https://api.system.mycompany.com
routing_cf_uaa_url: https://uaa.system.mycompany.com
routing_cf_app_domain_name: apps.mycompany.com
routing_cf_client_id: routing_api_client
routing-cf-client-secret: <<credhub get -n my-bosh/cf/uaa_clients_routing_api_client_secret>>
routing_cf_nats_internal_ips: [10.10.1.6,10.10.1.7,10.10.1.8]
routing_cf_nats_port: 4222
routing_cf_nats_username: nats
routing-cf-nats-password: <<credhub get -n my-bosh/cf/nats_password>>
```

You can try a helper script which might be able to use `bosh`, `cf`, and `credhub` CLIs to look up all the information:

```
./kubo-deployment/manifests/helper/cf-routing-vars.sh > cf-vars.yml
```

In the example `cf-vars.yml` above:

* the Cloud Foundry TCP router is available as hostname `tcp.apps.mycompany.com`, and route `tcp.apps.mycompany.com:8443` will be registered to route to the Kubernetes API running on all `master` instances of our deployment
* the Cloud Foundry internal NATS IPs are available via `bosh instances -d cf`
* extract the Credhub secrets and copy them into `cf-vars.yml`

    ```
    credhub get -n $BOSH_ENVIRONMENT/cf/uaa_clients_routing_api_client_secret --output-json | jq -r .value
    credhub get -n $BOSH_ENVIRONMENT/cf/nats_password --output-json | jq -r .value
    ```

NOTE: in future we can get rid of the `routing-cf-nats-*` variables and instead use the `nats` link from the `cf` deployment from the same BOSH. https://github.com/cloudfoundry-incubator/kubo-release/pull/134

NOTE: hopefully one day `cf` deployment will expose its URLs, admin credentials, and UAA clients via links and remove most of the other variables above. E.g. https://github.com/cloudfoundry/capi-release/pull/65

```
bosh deploy kubo-deployment/manifests/cfcr.yml \
  -o kubo-deployment/manifests/ops-files/cf-routing.yml \
  -l cf-vars.yml
```

We can now re-configure `kubectl` to use the new hostname and its matching certificate (rather than use the smelly `--insecure-skip-tls-verify` flag).

First, get the randomly generated Kubernetes API admin password from Credhub:

```
admin_password=$(bosh int <(credhub get -n "${BOSH_ENVIRONMENT}/${BOSH_DEPLOYMENT}/kubo-admin-password" --output-json) --path=/value)
```

Next, get your TCP hostname from your `cf-vars.yml` (e.g. `tcp.apps.mycompany.com`):

```
master_host=$(bosh int cf-vars.yml --path /kubernetes_master_host)
```

Then, store the root certificate in a temporary file:

```
tmp_ca_file="$(mktemp)"
bosh int <(credhub get -n "${BOSH_ENVIRONMENT}/${BOSH_DEPLOYMENT}/tls-kubernetes" --output-json) --path=/value/ca > "${tmp_ca_file}"
```

Finally, setup your local `kubectl` configuration:

```
cluster_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}"
user_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}-admin"
context_name="cfcr:${BOSH_ENVIRONMENT}:${BOSH_DEPLOYMENT}"

kubectl config set-cluster "${cluster_name}" \
  --server="https://${master_host}:8443" \
  --certificate-authority="${tmp_ca_file}" \
  --embed-certs=true
kubectl config set-credentials "${user_name}" --token="${admin_password}"
kubectl config set-context "${context_name}" --cluster="${cluster_name}" --user="${user_name}"
kubectl config use-context "${context_name}"
```

Confirm that the `:8443` TCP route and certificate for Kubernetes API are working:

```
kubectl get all
```

## Operator files

### BOSH options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`use-runtime-config-bosh-dns.yml`](use-runtime-config-bosh-dns.yml) | Delegate `bosh-dns` addon to BOSH runtime config | Apply this operator file if your BOSH environment has a runtime config that adds the `bosh-dns` job to all instances. By default, `cfcr.yml` will add `bosh-dns` to deployment instances. |
| [`vm-types.yml`](vm-types.yml) | Specify the `vm_type` for `master` and `worker` instances | By default, `master` and `worker` instances assume `vm_type: small` and `vm_type: small-highmem`, respectively (`vm_types` that are also assumed to exists by https://github.com/cloudfoundry/cf-deployment manifests). You may want to use bespoke `vm_types` so as to scale them, tag them, or apply unique `cloud_properties` independently of other deployments in the same BOSH environment. |

### Routing options

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`cf-routing.yml`](cf-routing.yml) | Combine CFCR with Cloud Foundry routing | Kube API and labeled services advertised to Cloud Foundry routing mesh. Kube API hostname is included in TLS certificates. |
| [`cf-routing-links.yml`](cf-routing-links.yml) | As above, using BOSH links | Simpler method of integration with Cloud Foundry running on same BOSH [in development] |
| **OpenStack** | | |
| [`iaas/openstack/master-static-ip.yml`](iaas/openstack/master-static-ip.yml) | Attach floating IP to Kube API | Assign allocated floating IP to `master` instance. IP included in TLS certificates. |
| **vSphere** | | |
| [`iaas/vsphere/master-static-ip.yml`](iaas/vsphere/master-static-ip.yml) | Assign static IP to Kube API | Assign static IP to `master` instance. IP included in TLS certificates. |
| **gcp** | | |
| [`iaas/gcp/add-service-key-master.yml`](iaas/gcp/add-service-key-master.yml) | Allow user to specify GCP key instead of service account |  |
| [`iaas/gcp/add-service-key-worker.yml`](iaas/gcp/add-service-key-worker.yml) | Allow user to specify GCP key instead of service account |  |
### Infrastructure

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| **AWS** | | |
| [`iaas/aws/cloud-provider.yml`](iaas/aws/cloud-provider.yml) | Enable Cloud Provider for AWS | Requires AWS Instance Profiles (not API keys) to grant Kubernetes access to AWS |
| [`iaas/aws/lb.yml`](iaas/aws/lb.yml) | Enable instance tagging for AWS |  |
| **OpenStack** | | |
| N/A | | |
| **GCP** | | |
| [`iaas/gcp/cloud-provider.yml`](iaas/gcp/cloud-provider.yml) | Enable Cloud Provider for OpenStack | - |
| **vSphere** | | |
| [`iaas/vsphere/cloud-provider.yml`](iaas/vsphere/cloud-provider.yml) | Enable Cloud Provider for vSphere | - |


### Proxy

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`add-http-proxy.yml`](add-http-proxy.yml) | Configure HTTP proxy for containers | Docker passes `http_proxy` environment variable to all containers |
| [`add-https-proxy.yml`](add-https-proxy.yml) | Configure HTTP proxy for containers | Docker passes `https_proxy` environment variable to all containers |
| [`add-no-proxy.yml`](add-no-proxy.yml) | Configure HTTP proxy for containers | Docker passes `no_proxy` environment variable to all containers |

### Kubernetes

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [`addons-spec.yml`](addons-spec.yml) | Addons to be deployed into the Kubernetes cluster | - |
