# Prerequisites for OpenStack Kubo deployment

1. An OpenStack environment running one of the following supported releases:

    - [Liberty](http://www.openstack.org/software/liberty)
    - [Mitaka](http://www.openstack.org/software/mitaka) (actively tested)
    - [Newton](http://www.openstack.org/software/newton)

1. The following OpenStack services:
   
    - [Identity](https://www.openstack.org/software/project-navigator#security)
    - [Compute](https://www.openstack.org/software/project-navigator#compute)
    - [Image](https://www.openstack.org/software/project-navigator#compute)
    - *(Recommended)* [OpenStack Networking](https://www.openstack.org/software/project-navigator#networking)
      provides network scaling and automated management functions.
      **Please note:** Nova networking is known to work, but is not actively
      tested, as it is deprecated.
    
1. An existing OpenStack project

1. The network should be configured to allow the following traffic:
    > **Note:** You can find the network settings in Compute > Access & Security, clicking the Manage Rules button for your security group.

    - Incoming TCP traffic for port 8443 for Kubernetes API 
    - Incoming TCP traffic for port range 30000-32765
    - Incoming UDP traffic for port 8285 within Kubo network
    - Incoming TCP traffic for ports 8844 and 25555 for operators
    - Outgoing TCP traffic for port 53 within Kubo network
    
1. *(Optional)* When using Cloud Foundry Routing with Kubo the following traffic
   should also be allowed:

    - Outgoing TCP traffic to CF routing API host and port
    - Outgoing TCP traffic to CF NATS hosts and port
    - Outgoing TCP traffic to CF router host and Kubo port
