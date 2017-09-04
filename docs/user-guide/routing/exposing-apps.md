# Exposing deployed applications directly

## Accessing an application on GCP and AWS with IaaS Load-Balancing

You can expose routes using the service type LoadBalancer for your Kubernetes deployments. See the [Kubernetes docs](https://kubernetes.io/docs/tutorials/kubernetes-basics/expose-intro/) for more details.

> **Note:** Any resources that are provisioned by Kubernetes will not be deleted by BOSH when you delete your Kubo deployment. You will need to manage these resources if they are not deleted by Kubernetes before the deployment is deleted.