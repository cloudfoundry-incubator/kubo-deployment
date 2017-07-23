# Kubo user guide

## Infrastructure paving
Kubo can be deployed on various infrastructure providers. The currently supported ones
are GCP and vSphere. vSphere installation is not yet covered in this guide.

Please follow the link for infrastructure paving on your platform:
- [Google Cloud Platform](platforms/gcp/paving.md)

## Deploying BOSH

- [Google Cloud Platform](platforms/gcp/install-bosh.md)
- [vSphere](platforms/vsphere/install-bosh.md)
- [OpenStack](platforms/openstack/install-bosh.md)

## Configure Kubo

### Create a Kubo environment

A Kubo environment is a set of configuration files used to deploy and update
both BOSH and Kubo.
You already have generated a KuBo environment if you completed the `Deploying BOSH` step above for 
your specific Iaas. In that case ignore this step and go to the next.

Run `./bin/generate_env_config <ENV_PATH> <ENV_NAME> gcp` to generate a Kubo
environment. The environment will be referred to as `KUBO_ENV` in this guide,
and will be located at `<ENV_PATH>/<ENV_NAME>`.

> Run `bin/generate_env_config --help` for more detailed information.

### Routing options

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

Once the infrastructure has been set up, a kubernetes cluster can be deployed by running a single line of code:

   ```bash
   bin/deploy_k8s <KUBO_ENV> <CLUSTER_NAME>
   ```

where `CLUSTER_NAME` is a unique name for the cluster. Run `bin/deploy_k8s --help` for more options on how to tell 
Bosh which release tarballs to use for the KuBo deployment (dev repo, internet, local, pre-uploaded to Bosh)

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

## Troubleshooting

See [troubleshooting section](troubleshooting.md) for solutions to most commonly encountered problems.
