# Customized Install (or Kubo the Hard Way)

## Dependencies
- [bosh-cli](https://bosh.io/docs/cli-v2.html) for interacting with BOSH. Version 2.0.1 and above. Please ensure the binary is installed as `bosh-cli` and not `bosh`.
- [credhub cli](https://github.com/pivotal-cf/credhub-cli/releases/tag/0.4.0) for interacting with CredHub. Version 0.4 only.
- [Ruby 2.3+](https://www.ruby-lang.org/en/downloads) required by the bosh-cli to deploy KuBOSH
- [make](https://www.gnu.org/software/make) required by the bosh-cli to deploy KuBOSH
- IaaS specific tools to create load balancers and firewall rules

## Infrastructure Setup

Kubo can be deployed with an IaaS load balancer that has a static external IP address and a forwarding rule to route traffic to the master nodes. If possible create a new subnet for Kubo to give it space and IP isolatuon. The following table specifies the needed routes and firewall rules

| Description                | Source                     | Destination              | Ports                                  |
|----------------------------|----------------------------|--------------------------|----------------------------------------|
| Access to KuBosh           | Your machine               | KuBosh                   | 2555/tcp, 6868/tcp, 8443/tcp, 8844/tcp |
| BOSH Management            | worker, master nodes       | KuBosh                   | 0-65535/tcp, 0-65535/udp               |
| Kubernetes API Endpoint    | Your machine, worker nodes | IaaS load balancer       | 8443/tcp                               |
| PowerDNS                   | worker, master nodes       | KuBosh                   | 53/tcp                                 |
| Kubernetes API routing     | IaaS load balancer         | master nodes             | 8443/tcp                               |

## Accessing the Kubo network

This guide assumes the machine you're executing commands from has proper access to the network Kubo will be deployed on. The easiest way to do this is to deploy a [jump box/bastion](https://en.wikipedia.org/wiki/Jump_server) on the Kubo subnet. You can then SSH into the machine and issue commands or use a tool like [sshuttle](https://github.com/apenwarr/sshuttle) to act as a VPN to your Kubo subnet.

## Configure and deploy KuBosh

Generate the configuration using templates for your IaaS with the following command:

```bash
bin/generate_env_config <path/to/generation/target/folder> <BOSH_NAME> <IAAS>
```

> Run `bin/generate_env_config --help` for more detailed information.

This will create a directory with the same name as the environment at the specified path, containing three files:
- `iaas` which contains IaaS name
- `director.yml` which contains public BOSH director, IaaS and network configurations. ([example](https://github.com/pivotal-cf-experimental/kubo-deployment/blob/master/ci/environments/gcp/director.yml))
- `director-secrets.yml` which contains sensitive configuration values, such as passwords and OAuth secrets

Edit these files to match your environment and generate a private key/service account with access to create/destroy VMs/disks then run the following command:

```bash
bin/deploy_bosh <BOSH_ENV> <private or service account key filename for BOSH to use for deployments> 
```

Credentials and SSL certificates for the environment will be generated and saved into the configuration path in a file called `creds.yml`. This file contains sensitive information and should not be stored in VCS. The file `state.json` contains state for the BOSH environment in the format of [bosh-init](https://bosh.io/docs/using-bosh-init.html). The `default` CA is generated and stored in CredHub.

Subsequent runs of `bin/bosh_deploy` will apply changes made to the configuration to an already existing KuBOSH installation, reusing the credentials stored in the `creds.yml`.

Another file that gets created during initial deployment is called `state.json`.

## Setup Cloud Config

Generate the Cloud Config and set it on your bosh director

```bash
bin/generate_cloud_config <BOSH_ENV> > <BOSH_ENV>/cloud-config.yml
# modify it as necessary
bosh-cli -e <BOSH_NAME> update-cloud-config <BOSH_ENV>/cloud-config.yml
```

## Generate manifest and deploy

Pick a deployment name and generate a manifest.

```bash
bin/generate_kubo_manifest <BOSH_ENV> <DEPLOYMENT_NAME> > <BOSH_ENV>/kubo-manifest.yml
```
The generation of the manifest can be customized in the following ways:

1. Variables in the manifest template can be substituted using an external file. Place a file named 
  `<DEPLOYMENT_NAME>-vars.yml` into the environment folder, and specify the variables as key-value
  pairs, e.g.:
  ```yaml
  super-secret-secret: SuperSecretPa$$phrase
  ```

1. Parts of the service manifest can be manipulated using 
  [go-patch](https://github.com/cppforlife/go-patch/blob/master/docs/examples.md) ops-files.
  To use this method, place a file named `<DEPLOYMENT_NAME>.yml` into the environment folder
  and fill it with go-patch instructions, e.g.:
  ```yaml
  - type: replace
    path: /releases/name=etcd
    value: 
      name: etcd
      version: 0.99.0
  ```

If needed, the generated manifest can be modified manually before being fed into `bosh-cli`:
```bash
bosh-cli -e <BOSH_NAME> -d <DEPLOYMENT_NAME> deploy <BOSH_ENV>/kubo-manifest.yml
```

## Accessing Kubo

Configure kubectl for your Kubo instance with the following command:

```bash
bin/set_kubeconfig <BOSH_ENV> <DEPLOYMENT_NAME>
```

You can now issue kubectl commands such as:
```bash
kubectl get pods --namespace=kube-system
kubectl get nodes
```

