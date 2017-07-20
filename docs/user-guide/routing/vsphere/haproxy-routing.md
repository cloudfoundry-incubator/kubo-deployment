# Load-balancing for vSphere using HAProxy

When deploying Kubo on vSphere, HAProxy can be used for external access to the Kubernetes Master node (for administration traffic), and the Kubernetes Workers (for application traffic). This is due to the fact that vSphere does not have first-party load-balancing support. In order to enable this, configure `director.yml` as follows:


Enable HAProxy routing:
```
routing_mode: proxy
```

Configure master IP address and port number:
```
kubernetes_master_host: 12.34.56.78
kubernetes_master_port: 443
```

Configure front-end and back-end ports for HAProxy TCP pass-through.
```
worker_haproxy_tcp_frontend_port: 1234
worker_haproxy_tcp_backend_port: 4321
```
*Note*: the current implementation of HAProxy routing is a single-port TCP pass-through. In order to route traffic to multiple Kubernetes services, consider using an Ingress Controller (https://github.com/kubernetes/ingress/tree/master/examples).
