# Customized Kubo installation 
## Setup Cloud Config

Generate the Cloud Config and set it on your bosh director

```bash
bin/generate_cloud_config <BOSH_ENV> > <BOSH_ENV>/cloud-config.yml
# modify cloud-config.yml as necessary
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

2. Parts of the service manifest can be manipulated using
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
