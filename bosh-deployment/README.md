# bosh-deployment

Requires BOSH CLI v2:

* [CLI v2](https://bosh.io/docs/cli-v2.html)
    * [`create-env` Dependencies](https://bosh.io/docs/cli-env-deps.html)
    * [Differences between CLI v2 vs v1](https://bosh.io/docs/cli-global-flags.html)
    * [Global Flags](https://bosh.io/docs/cli-global-flags.html)
    * [Environments](https://bosh.io/docs/cli-envs.html)
    * [Operations files](https://bosh.io/docs/cli-ops-files.html)
    * [Variable Interpolation](https://bosh.io/docs/cli-int.html)
    * [Tunneling](https://bosh.io/docs/cli-tunnel.html)

## Ops files

- `bosh.yml`: Base manifest that is meant to be used with different CPI configurations
- `[aws|azure|docker|gcp|openstack|softlayer|vcloud|vsphere|virtualbox]/cpi.yml`: CPI configuration
- `<cloud>/cloud-config.yml`: Simple cloud configs
- `jumpbox-user.yml`: Adds user `jumpbox` for SSH-ing into the Director (see [Jumpbox User](docs/jumpbox-user.md))
- `uaa.yml`: Deploys UAA and enables UAA user management in the Director
- `credhub.yml`: Deploys CredHub and enables CredHub integration in the Director
- `bosh-lite.yml`: Configures Director to use Garden CPI within the Director VM (see [BOSH Lite](docs/bosh-lite-on-vbox.md))
- `config-server.yml`: Deploys config-server (see `credhub.yml`)
- `syslog.yml`: Configures syslog to forward logs to some destination
- `local-dns.yml`: Enables Director DNS beta functionality
- `misc/proxy.yml`: Configure HTTP proxy for Director and CPI
- `runtime-configs/syslog.yml`: Runtime config to enable syslog forwarding

See [test.sh](test.sh) for example usage of different ops files.

## Sample installation instructions

* [BOSH Lite on VirtualBox](docs/bosh-lite-on-vbox.md)
* AWS (below)

```
$ git clone https://github.com/cloudfoundry/bosh-deployment ~/workspace/bosh-deployment

# Create a directory to keep Director deployment
$ mkdir -p ~/deployments/bosh-1

$ cd ~/deployments/bosh-1

# Deploy a Director -- ./creds.yml is generated automatically
$ bosh create-env ~/workspace/bosh-deployment/bosh.yml \
  --state ./state.json \
  -o ~/workspace/bosh-deployment/aws/cpi.yml \
  --vars-store ./creds.yml \
  -v access_key_id=... \
  -v secret_access_key=... \
  -v region=us-east-1 \
  -v az=us-east-1b \
  -v default_key_name=bosh \
  -v default_security_groups=[bosh] \
  -v subnet_id=subnet-... \
  -v director_name=bosh-1 \
  -v internal_cidr=10.0.0.0/24 \
  -v internal_gw=10.0.0.1 \
  -v internal_ip=10.0.0.6 \
  --var-file private_key=~/Downloads/bosh.pem

# Alias deployed Director
$ bosh -e 10.0.0.6 --ca-cert <(bosh int ./creds.yml --path /director_ssl/ca) alias-env bosh-1

# Log in to the Director
$ export BOSH_CLIENT=admin
$ export BOSH_CLIENT_SECRET=`bosh int ./creds.yml --path /admin_password`

# Update cloud config -- single az
$ bosh -e bosh-1 update-cloud-config ~/workspace/bosh-deployment/aws/cloud-config.yml \
  -v az=us-east-1b \
  -v subnet_id=subnet-... \
  -v internal_cidr=10.0.0.0/24 \
  -v internal_gw=10.0.0.1

# Upload specific stemcell
$ bosh -e bosh-1 upload-stemcell https://...

# Get a deployment running
$ git clone https://github.com/cppforlife/zookeeper-release ~/workspace/zookeeper-release
$ bosh -e bosh-1 -d zookeeper deploy ~/workspace/zookeeper-release/manifests/zookeeper.yml
```

To generate creds (without deploying anything) or just to check if your manifest builds:

```
$ bosh int ~/workspace/bosh-deployment/bosh.yml \
  --var-errs \
  -o ~/workspace/bosh-deployment/aws/cpi.yml \
  --vars-store ./creds.yml \
  -v access_key_id=... \
  -v secret_access_key=...
```

Please ensure you have security groups setup correctly. i.e:

```
Type                 Protocol Port Range  Source                     Purpose
SSH                  TCP      22          <IP you run bosh CLI from> SSH (if Registry is used)
Custom TCP Rule      TCP      6868        <IP you run bosh CLI from> Agent for bootstrapping
Custom TCP Rule      TCP      25555       <IP you run bosh CLI from> Director API
Custom TCP Rule      TCP      8443        <IP you run bosh CLI from> UAA API (if UAA is used)
SSH                  TCP      22          <((internal_cidr))>        BOSH SSH (optional)
Custom TCP Rule      TCP      4222        <((internal_cidr))>        NATS
Custom TCP Rule      TCP      25250       <((internal_cidr))>        Blobstore
Custom TCP Rule      TCP      25777       <((internal_cidr))>        Registry if enabled
```
