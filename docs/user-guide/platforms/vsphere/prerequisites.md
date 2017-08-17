# Prerequisites for vSphere deployment

The following details are needed in order to deploy Kubo on vSphere.

## General

On vSphere, Bosh needs a user account with a particular set of privileges. This guide will
refer to this account as the _bosh user_. Please refer to
[vSphere Documentation Center](http://pubs.vmware.com/vsphere-65/index.jsp?topic=%2Fcom.vmware.vsphere.security.doc%2FGUID-18071E9A-EED1-4968-8D51-E0B4F526FDA3.html&resultof=%22%43%72%65%61%74%65%22%20%22%63%72%65%61%74%22%20%22%43%75%73%74%6f%6d%22%20%22%63%75%73%74%6f%6d%22%20%22%52%6f%6c%65%22%20%22%72%6f%6c%65%22%20)
for more details on user accounts and roles. 

Please make sure to have all the following details:

1. A vSphere role associated with the _bosh user_ should grant the following privileges:    
    
    - Data Store
        - Allocate space
        - Browse datastore
        - Low level file operations
        - Remove file
        - Update virtual machine files
        - Update virtual machine metadata
    - Folder
    - Global
        - Manage custom attributes
        - Set custom attribute
    - Host
        - Inventory
            - Modify cluster
        - Local operations
    - Network
    - Resource
        - Assign virtual machine to resource pool
        - Migrate powered off virtual machine
        - Migrate powered on virtual machine
    - Virtual Machine
        - Configuration
        - Guest Operations
        - Interaction
        - Inventory
        - Provisioning
        - Service configuration
        - Snapshot management
    - vApp
    - vCenter Inventory Service

1. Ensure you have the following vSphere configuration details. These values will be required later during the deployment of BOSH and Kubo.
    - vCenter IP address
    - Username and password for the _bosh user_
    - A vSphere datacenter name
    - A name of an existing cluster in the datacenter
    - A name for an existing datastore in the same datacenter
    - A name for an existing resource pool in the cluster
    - A name for a Templates folder
    - A name for a VMs folder
    - A name for a disks folder

The folders mentioned above will be created during deployment if they do not exist at that time.

## Kubernetes persistence role

**(Optional)** When Kubernetes applications have to access persistent volumes, it is recommended to create a
separate user account with a tighter set of privileges using the guidelines below. This user account will be 
referred to as the _persistence user_ in this guide.

1. Username and password for the _persistence user_
1. The role associated with the _persistence user_ should grant the following privileges:
    
    - Datastore
        - Allocate space
        - Low level file Operations
    - Virtual Machine
        - Configuration
            - Add existing disk
            - Add or remove device
            - Remove disk

1. _(Optional)_ VSAN policy based volume provisioning feature in Kubernetes requires the 
following additional privileges:
    
    - Network
        - Assign network
    - Virtual machine
        - Configuration
            - Add new disk
    - Virtual Machine
        - Inventory
            - Create new
    - Resource
        - Assign virtual machine to resource pool
