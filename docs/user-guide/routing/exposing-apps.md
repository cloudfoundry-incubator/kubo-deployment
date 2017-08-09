# Exposing deployed applications directly

## Accessing an application on GCP with IaaS Load-Balancing

You can expose routes using the service type LoadBalancer for your Kubernetes deployments. See the [Kubernetes docs](https://kubernetes.io/docs/tutorials/kubernetes-basics/expose-intro/) for more details.

> **Note:** Any resources that are provisioned by Kubernetes will not be deleted by BOSH when you delete your Kubo deployment. You will need to manage these resources if they are not deleted by Kubernetes before the deployment is deleted.

## Accessing an application on AWS with IaaS Load-Balancing

An additional load balancer is provisioned using Terraform during the setup in our guide.
You can access the service using the external IP address of the  `kubo-apps` load
balancer and the `NodePort` of your service.

1. Add a listener on the `kubo-workers` load balancer for the `NodePort` of your service

1. Find the DNS name of your load balancer in the Description section

1. Access your service at `<DNS name of load balancer>:<NodePort>`
