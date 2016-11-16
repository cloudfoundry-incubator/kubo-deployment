# bosh-deployment

- Requires new [BOSH CLI v0.0.112+](https://github.com/cloudfoundry/bosh-cli)

AWS:

```
$ mkdir bosh-1 && cd bosh-1

$ cp ~/workspace/bosh-deployment/bosh.yml ./manifest.yml

$ bosh create-env ./manifest.yml \
  --ops-file ~/workspace/bosh-deployment/use-aws.yml \
  --vars-store ./creds.yml \
  -v access_key_id=... \
  -v secret_access_key=... \
  -v region=us-east-1 \
  -v az=us-east-1b \
  -v default_key_name=$(dirname $PWD) \
  -v default_security_groups=[bosh] \
  -v subnet_id=subnet-... \
  -v director_name=my-bosh \
  -v internal_cidr=10.10.0.0/24 \
  -v internal_gw=10.10.0.1 \
  -v internal_ip=10.10.0.6 \
  -v private_key=...
```
