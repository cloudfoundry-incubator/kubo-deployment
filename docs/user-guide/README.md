# Kubo user guide

Kubo can be deployed on various infrastructure providers. The currently supported ones
are GCP, vSphere and OpenStack.

## Prerequisites

- [GCP](platforms/gcp/prerequisites.md)
- [vSphere](platforms/vsphere/prerequisites.md)

## Infrastructure paving

Some platforms allow automatic infrastructure paving in order to prepare
an environment for a Kubo deployment. Please follow the link for infrastructure 
paving on your platform:

- [Google Cloud Platform](platforms/gcp/paving.md)
- [Manual setup](manual-paving.md)

## Deploying BOSH

- [Google Cloud Platform](platforms/gcp/install-bosh.md)
- [vSphere](platforms/vsphere/install-bosh.md)
- [OpenStack](platforms/openstack/install-bosh.md)
- [Manual deployment](bosh-customized-installation.md)

## Configure Kubo

### Create a Kubo environment

A Kubo environment is a set of configuration files used to deploy and update both BOSH and Kubo. If you followed the [Deploying BOSH](#deploying-bosh) step above for 
your specific IaaS, ignore this step and go to [Routing options](#routing_options).

Otherwise, run `./bin/generate_env_config <ENV_PATH> <ENV_NAME> <PLATFORM_NAME>`
to generate a Kubo environment. The environment will be referred to as `KUBO_ENV` 
in this guide, and will be located at `<ENV_PATH>/<ENV_NAME>`.

> Run `bin/generate_env_config --help` for more detailed information.

### <a name="routing_options">Routing options</a>

Kubo can leverage different routing options to improve security and high
availability of the cluster. Please configure the kubo environment according
to one of the options below:

- [IaaS Load-Balancing](routing/gcp/load-balancing.md)
- [CF Routing](routing/cf.md)
- [HAProxy Routing for OpenStack](routing/openstack/haproxy-routing.md)
- [HAProxy Routing for vSphere](routing/vsphere/haproxy-routing.md)

### Basic configuration

To view an adjust the available configuration options, please edit the `<KUBO_ENV>/director.yml` file that
was created as part of the [Kubo environment creation step](#create-a-kubo-environment).

### Proxy settings

The following variables can be configured in the director.yml to allow Docker have proxy access:

```yaml
http_proxy: # e.g. http://my.proxy.local:73636
https_proxy: # e.g. https://secure.proxy.local:5566
no_proxy: # e.g. '1.2.3.4,2.3.4.5'
```

## Deploying Kubo

> Get latest version of kubo-deployment before deployment. Deploying K8s using public tarballs using old version of kubo-deployment
> might result in an error.

Once the infrastructure has been set up, a Kubernetes cluster can be deployed by running a single line of code:

   ```bash
   cd /share/kubo-deployment
   bin/deploy_k8s <KUBO_ENV> <MY_CUSTOM_CLUSTER_NAME>
   ```

where `<KUBO_ENV>` is located at `<ENV_PATH>/<ENV_NAME>` and where `CLUSTER_NAME` is a unique name for the cluster. 
Run `bin/deploy_k8s --help` for more options on how to tell Bosh which release tarballs to use for the KuBo deployment:
manually built from repo(`dev`), precompiled from internet(`public`),
manually downloaded to specific location(`local`), pre-uploaded to Bosh(`skip`).

### Customized deployment

The `kubo-deployment` provides a number of ways to customize the kubo settings. Please follow the
[custom install guide](customized-kubo-installation.md) if you need to change the default behaviour.

## Accessing Kubernetes


### Operator access
Once the cluster is deployed, setup `kubectl` and access your new Kubernetes cluster

   ```bash
   bin/set_kubeconfig <KUBO_ENV> <CLUSTER_NAME> 
   ```

To verify that the settings have been applied correctly, run the following command:

   ```bash
   kubectl get pods --namespace=kube-system
   ```

### Enabling application access

Different routing modes provide different ways of exposing applications run by the Kubernetes cluster:

- [IaaS routing](./routing/exposing-apps.md)
- [CF routing](./routing/cf-apps.md)

### Persistence

Kubo clusters currently support the following Kubernetes Volume types:
- emptyDir
- hostPath
- gcePersistentDisk
- VsphereVolume

To use storage in the Kubo clusters the `cloud-provider` job must be configured on the master and worker instances. See the [cloud-provider spec](https://github.com/cloudfoundry-incubator/kubo-release/blob/master/jobs/cloud-provider/spec) for details on the properties that are needed for each cloud-provider type. 

For documentation on configuring Kubernetes to access storage for your cloud-provider type see - https://kubernetes.io/docs/concepts/storage/persistent-volumes/

> **Note:** Any resources that are provisioned by Kubernetes will not be deleted by BOSH when you delete your Kubo deployment. You will need to manage these resources if they are not deleted by Kubernetes before the deployment is deleted.

## Troubleshooting

See [troubleshooting section](troubleshooting.md) for solutions to most commonly encountered problems.
