# Kubo User guide

1. [Prerequisites](#prerequisites)
1. [Infrastructure Paving](#infrastructure-paving)
1. [Deploying BOSH](#deploying-bosh)
1. [Configure Kubo](#configure-kubo)
1. [Deploying Kubo](#deploying-kubo)
1. [Accessing Kubernetes](#accessing-kubernetes)
1. [Destroying the Cluster](#destroying-the-cluster)
1. [Destroying the BOSH Environment](#destroying-the-bosh-environment)
1. [Troubleshooting](#troubleshooting)

Kubo can be deployed on various infrastructure providers. The currently supported ones
are GCP, vSphere, AWS and OpenStack.

## Prerequisites

- [GCP](platforms/gcp/prerequisites.md)
- [vSphere](platforms/vsphere/prerequisites.md)
- [AWS](platforms/aws/prerequisites.md)
- [OpenStack](platforms/openstack/prerequisites.md)

## Infrastructure Paving

Some platforms allow automatic infrastructure paving in order to prepare
an environment for a Kubo deployment. Please follow the link for infrastructure
paving on your platform:

- [GCP](platforms/gcp/paving.md)
- [AWS](platforms/aws/paving.md)

## Deploying BOSH

- [GCP](platforms/gcp/install-bosh.md)
- [vSphere](platforms/vsphere/install-bosh.md)
- [AWS](platforms/aws/install-bosh.md)
- [OpenStack](platforms/openstack/install-bosh.md)

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

- [IaaS Load-Balancing on GCP](routing/gcp/load-balancing.md)
- [IaaS Load-Balancing on AWS](routing/aws/load-balancing.md)
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

Once the infrastructure has been set up, a Kubernetes cluster can be deployed by running `deploy_k8s`. This command will download all the packages necessary to deploy Kubernetes, and then bring up the master, worker, and etcd nodes as a managed cluster:

```bash
# From the kubo-deployment directory:
bin/deploy_k8s <KUBO_ENV> <MY_CLUSTER_NAME>
```

`KUBO_ENV` is located at `<ENV_PATH>/<ENV_NAME>` and `MY_CLUSTER_NAME` is a unique name for the cluster. Run `bin/deploy_k8s --help` for more options on how to tell BOSH which release tarballs to use for the Kubo deployment:

* manually built from the local repo (`dev`)
* precompiled from internet (`public`) - This is the default option
* manually downloaded to a specific location (`local`)
* pre-uploaded to BOSH (`skip`)

### Customized deployment

The `kubo-deployment` provides a number of ways to customize Kubo settings. Please follow the [custom install guide](customized-kubo-installation.md) if you need to change the default behaviour.

## Accessing Kubernetes

### Operator access

Once the cluster is deployed, setup `kubectl` and access your new Kubernetes cluster

```bash
bin/set_kubeconfig <KUBO_ENV> <MY_CLUSTER_NAME>
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
- AWSElasticBlockStore

To use storage in the Kubo clusters the `cloud-provider` job must be configured on the master and worker instances. See the [cloud-provider spec](https://github.com/cloudfoundry-incubator/kubo-release/blob/master/jobs/cloud-provider/spec) for details on the properties that are needed for each cloud-provider type.

For documentation on configuring Kubernetes to access storage for your cloud-provider type see - https://kubernetes.io/docs/concepts/storage/persistent-volumes/

> **Note:** Any resources that are provisioned by Kubernetes will not be deleted by BOSH when you delete your Kubo deployment. You will need to manage these resources if they are not deleted by Kubernetes before the deployment is deleted.

## Destroying the Cluster

Use the BOSH CLI if you want to destroy the cluster:

```bash
bosh-cli -e <KUBO_ENV> login
bosh-cli -e <KUBO_ENV_NAME> -d <MY_CLUSTER_NAME> delete-deployment
```

Your username is admin and your password is the `admin_password` field in `<KUBO_ENV>/creds.yml`.

`KUBO_ENV_NAME` was set up in the Install BOSH step.

## Destroying the BOSH Environment

To destroy your BOSH environment, follow the guide for your specific platform:

* [GCP](platforms/gcp/destroy-bosh.md)

## Troubleshooting

See [troubleshooting section](troubleshooting.md) for solutions to most commonly encountered problems.
