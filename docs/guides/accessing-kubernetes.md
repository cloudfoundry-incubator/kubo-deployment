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
  
  
    curl 10.244.4.11:31043
  ```
  
  This should display the standard nginx welcome page.

## Accessing Kubernetes dashboard

1. Get the NodePort for the Kubernetes dashboard service

   ```
   kubectl describe service kubernetes-dashboard --namespace=kube-system
   ```

1. Get the IP address of one of the Kubernetes worker nodes

   ```
   bosh-cli -e kube vms
   ```

1. Setup [sshuttle](https://github.com/apenwarr/sshuttle) from your local machine to your KuBOSH Director

   ```
   sshuttle -r <local machine IP address> <KuBOSH Subnet CIDR>
   ```

1. View the Kubernetes dashboard from your browser at `<worker node IP>:<NodePort>`
