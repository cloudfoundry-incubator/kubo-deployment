# Exposing deployed applications directly

The only way to expose Kubernetes applications running in Kubo is to use a 
[`NodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport). 

We do not currently support the type [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#type-loadbalancer),
but we plan to soon with Github issue [#47](https://github.com/cloudfoundry-incubator/kubo-release/issues/47)
in the [kubo-release](https://github.com/cloudfoundry-incubator/kubo-release) 
repository. 

## Accessing an application on GCP with IaaS Load-Balancing

An additional load balancer is provisioned using Terraform during the setup in our guide. 
You can access the service using the external IP address of the  `kubo-workers` load 
balancer and the `NodePort` of your service.

1. Access your service from your browser at `<IP address of the load balancer>:<NodePort>`
