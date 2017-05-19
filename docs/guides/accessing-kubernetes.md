# Kubernetes Usage Examples

## Accessing the Kubernetes cluster

1. Setup kubectl
   ```
   bin/set_kubeconfig ${state_dir} kube
   ```

1. Access your Kubernetes cluster
   ```
   kubectl get pods --namespace=kube-system
   ```

## Deploy a Kubernetes workload

1. Once you can access your cluster, you should be able to deploy a simple workload. We'll try nginx as an example:
   ```
   kubectl create -f ci/specs/nginx.yml
   ```

1. This should create 3 nginx pods:

   ```
   kubectl get pods -o wide
   
     NAME                     READY     STATUS    RESTARTS   AGE       IP            NODE
     nginx-2793416118-67icb   1/1       Running   0          36s       10.200.40.2   10.244.4.11
     nginx-2793416118-6q6x7   1/1       Running   0          36s       10.200.26.2   10.244.4.20
     nginx-2793416118-zyicv   1/1       Running   0          36s       10.200.62.2   10.244.4.10
   ```

1. Get the IP address of one of the Kubernetes worker nodes

   ```
   bosh-cli -e kube vms

     Instance                                     Process State  AZ  IPs        VM CID                                   VM    Type
     etcd/bdaee053-f5e3-4f9d-b27e-02072595fc70    running        z1  10.0.1.2   vm-936a223a-5be0-487c-6e8c-b75ae1cb0015  common
     etcd/c69ff1cd-050a-4201-bc73-b501b14fc71a    running        z1  10.0.1.6   vm-bd8f16d1-a582-4288-680f-50aa37c9c6f5  common
     etcd/daef49a9-16e4-48c6-a362-85c560cf3d77    running        z1  10.0.1.12  vm-b7fa754c-906c-4c23-5850-312680e31abb  common
     master/6b9a9502-da1c-4ae0-aba0-2eb06974be78  failing        z1  10.0.1.8   vm-91c08282-5c53-4523-6d5d-f94d7ac0b837  common
     master/b9c3708b-8a0c-441d-9f83-4f4d19ef5234  running        z1  10.0.1.10  vm-3d7f9920-7bff-4463-5606-8aa1de747a95  common
     proxy/b4820ea6-4fe8-4697-a87f-3286cd507f90   running        z1  10.0.1.5   vm-034f6216-7f06-4906-6aeb-039544f0f0c7  common
     worker/1d878f78-0a2c-458d-b611-801a706fb74a  running        z1  10.0.1.11  vm-16ba24eb-7cc6-46c4-73ff-4033d80eb3b2  worker
     worker/6897a986-14ae-4ea3-8849-4c930f903a7b  running        z1  10.0.1.9   vm-8077f674-1292-4f22-52c8-3b875c058fca  worker
     worker/ae232a39-0930-4482-9213-c6ea5e1ddab1  running        z1  10.0.1.7   vm-04bcebfa-b2a8-4f40-6363-3762030ea52c  worker
   ```

1. Test that nginx is really running by obtaining the service's NodePort, and curling that port on any of your worker nodes:

   ```
   kubectl describe svc nginx
   
     Name:            nginx
     Namespace:        default
     Labels:            name=nginx
     Selector:        app=nginx
     Type:            NodePort
     IP:            10.100.200.199
     Port:            <unset>    80/TCP
     NodePort:        <unset>    31043/TCP
     Endpoints:        10.200.26.2:80,10.200.40.2:80,10.200.62.2:80
     Session Affinity:    None
     No events.
   
   
   curl <worker node IP>:<NodePort>
   ```
   
   This should display the standard nginx welcome page.

## Accessing Kubernetes services

To expose services running in your Kubernetes cluster, use the service type `NodePort` when deploying your service to Kubernetes. We do not currently support the type `LoadBalancer`, but we plan to soon with Github issue [#47](https://github.com/pivotal-cf-experimental/kubo-release/issues/47) in the [kubo-release](https://github.com/pivotal-cf-experimental/kubo-release) repository. Until this issue is resolved, an additional load balancer is provisioned using Terraform during the setup in our guide. If your service is exposed with a NodePort, you can access the service using the external IP address of the kubo-workers load balancer and the node port of your service.

### Example: Accessing the Kubernetes dashboard on GCP
   
1. Find the IP address of your worker load balancer

   ```
   gcloud compute addresses list | grep kubo-workers

     kubo-workers    us-west1  XX.XXX.X.XXX     IN_USE
   ```

1. Find the Node Port of the kubernetes-dashboard service

   ```
   kubectl describe service kubernetes-dashboard --namespace kube-system  | grep NodePort

     Type:                   NodePort
     NodePort:               <unset> 31000/TCP
   ```

1. Access your service from your browser at `<IP address of load balancer>:<NodePort>`
