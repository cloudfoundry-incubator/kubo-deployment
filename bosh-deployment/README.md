# bosh-deployment

* [Create an environment](https://bosh.io/docs/init.html)
    * [On Local machine (BOSH Lite)](https://bosh.io/docs/bosh-lite.html)
    * [On AWS](https://bosh.io/docs/init-aws.html)
    * [On Azure](https://bosh.io/docs/init-azure.html)
    * [On OpenStack](https://bosh.io/docs/init-openstack.html)
    * [On vSphere](https://bosh.io/docs/init-vsphere.html)
    * [On vCloud](https://bosh.io/docs/init-vcloud.html)
    * [On SoftLayer](https://bosh.io/docs/init-softlayer.html)
    * [On Google Compute Platform](https://bosh.io/docs/init-google.html)

* Access your BOSH director
    * Through a VPN
        * [`bosh create-env`, OpenVPN option](https://github.com/dpb587/openvpn-bosh-release)
    * Through a jumpbox
        * [`bosh create-env` option](https://github.com/cppforlife/jumpbox-deployment)
    * [Expose Director on a Public IP](https://bosh.io/docs/init-external-ip.html) (not recommended)

* [CLI v2](https://bosh.io/docs/cli-v2.html)
    * [`create-env` Dependencies](https://bosh.io/docs/cli-v2-install/#additional-dependencies)
    * [Differences between CLI v2 vs v1](https://bosh.io/docs/cli-v2-diff.html)
    * [Global Flags](https://bosh.io/docs/cli-global-flags.html)
    * [Environments](https://bosh.io/docs/cli-envs.html)
    * [Operations files](https://bosh.io/docs/cli-ops-files.html)
    * [Variable Interpolation](https://bosh.io/docs/cli-int.html)
    * [Tunneling](https://bosh.io/docs/cli-tunnel.html)

## Ops files

- `bosh.yml`: Base manifest that is meant to be used with different CPI configurations
- `[aws|azure|docker|gcp|openstack|softlayer|vcloud|vsphere|virtualbox]/cpi.yml`: CPI configuration
- `[aws|azure|docker|gcp|openstack|softlayer|vcloud|vsphere|virtualbox]/cloud-config.yml`: Simple cloud configs
- `jumpbox-user.yml`: Adds user `jumpbox` for SSH-ing into the Director (see [Jumpbox User](docs/jumpbox-user.md))
- `uaa.yml`: Deploys UAA and enables UAA user management in the Director
- `credhub.yml`: Deploys CredHub and enables CredHub integration in the Director
- `bosh-lite.yml`: Configures Director to use Garden CPI within the Director VM (see [BOSH Lite](docs/bosh-lite-on-vbox.md))
- `syslog.yml`: Configures syslog to forward logs to some destination
- `local-dns.yml`: Enables Director DNS beta functionality
- `misc/config-server.yml`: Deploys config-server (see `credhub.yml`)
- `misc/proxy.yml`: Configure HTTP proxy for Director and CPI
- `runtime-configs/syslog.yml`: Runtime config to enable syslog forwarding

See [tests/run-checks.sh](tests/run-checks.sh) for example usage of different ops files.

## Security Groups

Please ensure you have security groups setup correctly. i.e:

```
Type                 Protocol Port Range  Source                     Purpose
SSH                  TCP      22          <IP you run bosh CLI from> SSH (if Registry is used)
Custom TCP Rule      TCP      6868        <IP you run bosh CLI from> Agent for bootstrapping
Custom TCP Rule      TCP      25555       <IP you run bosh CLI from> Director API
Custom TCP Rule      TCP      8443        <IP you run bosh CLI from> UAA API (if UAA is used)
Custom TCP Rule      TCP      8844        <IP you run bosh CLI from> CredHub API (if CredHub is used)
SSH                  TCP      22          <((internal_cidr))>        BOSH SSH (optional)
Custom TCP Rule      TCP      4222        <((internal_cidr))>        NATS
Custom TCP Rule      TCP      25250       <((internal_cidr))>        Blobstore
Custom TCP Rule      TCP      25777       <((internal_cidr))>        Registry if enabled
```
