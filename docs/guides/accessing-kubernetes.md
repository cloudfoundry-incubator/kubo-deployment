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

If you deployed Kubo using the Cloud Foundry routers, you can expose routes to your services in the following way:

### Creating TCP Routes
1. Add a label to your service where the label is named `tcp-route-sync` and the value of the label is the frontend port that you want to expose your application on
   ```
   kubectl label services <your service name> tcp-route-sync=<frontend port>
   ```

1. Access your service from your browser at `<Cloud Foundry tcp url>:<frontend port>`

   > **Note:** It may take up to 60 seconds for the route to be created

### Creating HTTP Routes
1. Add a label to your service where the label is named `http-route-sync` and the value of the label is the name of the route that you want to create for your application
   ```
   kubectl label services <your service name> http-route-sync=<route name>
   ```
   
1. Access your service from your browser at `<route name>.<Cloud Foundry apps domain>`
   
   > **Note:** It may take up to 60 seconds for the route to be created
   
### Example: Accessing the Kubernetes dashboard
   
1. Expose an HTTP route for the dashboard service
   ```
   kubectl label services kubernetes-dashboard http-route-sync=dashboard --namespace=kube-system
   ```
   
1. View the Kubernetes dashboard from your browser at `dashboard.<Cloud Foundry apps domain>`
