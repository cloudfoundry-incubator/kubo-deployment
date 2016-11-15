# bosh-deployment

- Requires new [BOSH CLI](https://github.com/cloudfoundry/bosh-cli)

`bosh create-env gcp.yml -l path/to/vars/file.yml --state `

common:

internal_cidr
internal_gw
internal_ip
nats_password
postgres_password
blobstore_director_password
agent_director_password
admin_password
hm_password
director_name
director_ssl_private_key
director_ssl_certificate
director_ssl_ca

gcp: zone, service_account, network, subnetwork, project_id
aws:
vars:
  - key: az
    title: "something helpful"
  - subnet_id, registry_password, access_key_id, secret_access_key, region, private_key
