# Deploying Kubernetes using BOSH

## Easy way

1. Run this command:
   ```bash
   bin/deploy_k8s <BOSH_ENV> <DEPLOYMENT_NAME> public
   ```

## Custom installation

1. Generate and apply [cloud config](https://bosh.io/docs/cloud-config.html)
    ```
    bin/generate_cloud_config <BOSH_ENV> > <BOSH_ENV>/cloud-config.yml
    ```
    Modify `<BOSH_ENV>/cloud-config.yml` as you wish and apply it. 

1. Generate [deployment manifest](https://bosh.io/docs/manifest-v2.html)
    ```
    bin/generate_kubo_manifest <BOSH_ENV> <DEPLOYMENT_NAME> > <BOSH_ENV>/manifest.yml
    ```
    Modify manifest using [operation files](https://bosh.io/docs/cli-ops-files.html).

1. Upload [releases](https://bosh.io/docs/release.html) and [stemcell](https://bosh.io/docs/stemcell.html).
    Check the version and required releases in the manifest.

1. Deploy result manifest