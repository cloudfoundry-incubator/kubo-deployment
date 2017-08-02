# Prerequisites for OpenStack Kubo deployment

1. An OpenStack environment running one of the following supported releases:

    - [Liberty](http://www.openstack.org/software/liberty)
    - [Mitaka](http://www.openstack.org/software/mitaka) (actively tested)
    - [Newton](http://www.openstack.org/software/newton)

1. The following OpenStack services:
   
    - [Identity](http://www.openstack.org/software/openstack-shared-services/): BOSH authenticates credentials and 
      retrieves the endpoint URLs for other OpenStack services.
    - [Compute](http://www.openstack.org/software/openstack-compute/): BOSH boots new VMs, assigns floating IPs to VMs, 
      and creates and attaches volumes to VMs.
    - [Image](http://www.openstack.org/software/openstack-shared-services/): BOSH stores stemcells using the Image service.
    - *(Recommended)* [OpenStack Networking](https://www.openstack.org/software/): Provides network scaling and automated management
      functions that are useful when deploying complex distributed systems. **Note:** Nova networking is known to work,
      but is not actively tested, as it is deprecated.
    
1. An existing OpenStack project

1. The network should be configured to allow the following traffic:

    - Incoming TCP traffic for port 8443 for Kubernetes API 
    - Incoming TCP traffic for port range 30000-32765
    - Incoming UDP traffic for port 8285 within Kubo network
    - Incoming TCP traffic for ports 8844 and 25555 for operators
    - Outgoing TCP traffic for port 53 within Kubo network
    
1. *(Optional)* When using CF routing with Kubo the following traffic should also be allowed:

    - Outgoing TCP traffic to CF routing API host and port
    - Outgoing TCP traffic to CF NATS hosts and port
    - Outgoing TCP traffic to CF router host and Kubo port