# Destroy the BOSH Environment on GCP

Execute `destroy_bosh` to destroy all the VMs that make up the BOSH environment.

```bash
bin/destroy_bosh <KUBO_ENV> ~/terraform.key.json
```

`terraform.key.json` is the GCP service account key used in the [Install BOSH](install-bosh.md) step.

## Destroy the Infrastructure Paved by Terraform

Use Terraform to destroy the resources created in the [Paving](paving.md) step. You'll need the same environmnent variables as in that step. When you run `destroy`, make sure you're in the same directory you executed `apply` in the Paving step and that the `terraform.tfstate` is there.

```bash
cd ~/kubo-deployment/docs/user-guide/platforms/gcp
docker run -i -t \
    -e CHECKPOINT_DISABLE=1 \
    -e "GOOGLE_CREDENTIALS=${GOOGLE_CREDENTIALS}" \
    -v $(pwd):/$(basename $(pwd)) \
    -w /$(basename $(pwd)) \
    hashicorp/terraform:light destroy \
    -var service_account_email=${service_account_email} \
    -var projectid=${project_id} \
    -var network=${network} \
    -var region=${region} \
    -var prefix=${prefix} \
    -var zone=${zone} \
    -var subnet_ip_prefix=${subnet_ip_prefix}
```
