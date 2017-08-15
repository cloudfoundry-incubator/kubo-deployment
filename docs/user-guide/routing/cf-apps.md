# Enabling application access via CF routers
[Cloud Foundry routers](https://docs.cloudfoundry.org/devguide/deploy-apps/routes-domains.html#http-vs-tcp-routes)
can be used to expose both TCP and HTTP level routes to the Kubernetes services.

## Preconditions
1. K8s service must be exposed using a single `NodePort`

## Creating TCP Routes
1. Add a label to your service where the label is named `tcp-route-sync` and the value of the label is the frontend port that you want to expose your application on
   ```
   kubectl label services <your service name> tcp-route-sync=<frontend port>
   ```
   > **Note:** The frontend port must be a port within the range you configured when enabling TCP Routing in CF

1. Access your service from your browser at `<Cloud Foundry tcp url>:<frontend port>`

   > **Note:** It may take up to 60 seconds for the route to be created

## Creating HTTP Routes
1. Add a label to your service where the label is named `http-route-sync` and the value of the label is the name of the route that you want to create for your application
   ```
   kubectl label services <your service name> http-route-sync=<route name>
   ```

1. Access your service from your browser at `<route name>.<Cloud Foundry apps domain>`

   > **Note:** It may take up to 60 seconds for the route to be created
