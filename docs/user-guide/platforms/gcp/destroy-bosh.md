# Destroy the BOSH environment on GCP

Execute `destroy_bosh` to destroy all the VMs that make up the BOSH environment.

```bash
bin/destroy_bosh <KUBO_ENV> ~/terraform.key.json
```

`terraform.key.json` is the GCP used in the [Install BOSH](install-bosh.md) step.

## Destroy the infrastructure paved by terraform

Use Terraform to destroy the resources created in the [Paving](paving.md) step. You'll need the same environmnent variables as in that step. Also, you'll need to change to the directory where your `terraform.tfstate` was created (most likely the same directory you executed `apply` in the Paving step).

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
