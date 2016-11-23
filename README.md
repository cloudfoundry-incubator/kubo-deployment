# bosh-deployment

- Requires new [BOSH CLI v0.0.120+](https://github.com/cloudfoundry/bosh-cli)

AWS:

```
$ cd ~/workspace

$ git clone https://github.com/cloudfoundry/bosh-deployment

# Create a directory to keep Director deployment
$ mkdir bosh-1 && cd bosh-1

# Deploy a Director
$ bosh create-env ~/workspace/bosh-deployment/bosh.yml \
  --state ./state.json \
  --ops-file ~/workspace/bosh-deployment/use-aws.yml \
  --vars-store ./creds.yml \
  -v access_key_id=... \
  -v secret_access_key=... \
  -v region=us-east-1 \
  -v az=us-east-1b \
  -v default_key_name=bosh \
  -v default_security_groups=[bosh] \
  -v subnet_id=subnet-... \
  -v director_name=$(dirname $PWD) \
  -v internal_cidr=10.10.0.0/24 \
  -v internal_gw=10.10.0.1 \
  -v internal_ip=10.10.0.6 \
  -v private_key=...

# Alias deployed Director
$ bosh alias-env bosh-1 \
  -e `bosh int ./creds.yml --path /internal_ip` \
  --ca-cert <(bosh int ./creds.yml --path /director_ssl/ca)

# Log in to the Director
$ bosh -e bosh-1 log-in
```

To generate creds (without deploying anything) or just to check if your manifest builds:

```
$ bosh int ~/workspace/bosh-deployment/bosh.yml \
  --var-errs \
  --ops-file ~/workspace/bosh-deployment/use-aws.yml \
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
