# Destroy the BOSH environment on GCP

Execute `destroy_bosh` to destroy all the VMs that make up the BOSH environment.

```bash
bin/destroy_bosh <KUBO_ENV> ~/terraform.key.json
```

`terraform.key.json` is the GCP used in the [Install BOSH](install-bosh.md) step.
