# Prerequisites for vSphere deployment

The following details are needed to deploy Kubo on vSphere:

1. A vSphere role with the following privileges:    
    
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
    - Inventory Service
        - vSphere Tagging
            - Create vSphere Tag
            - Delete vSphere Tag
            - Edit vSphere Tag
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

1. vCenter IP address
1. credentials for a user account associated with the role outlined above
1. a vSphere datacenter name
1. a name of the cluster in the datacenter
1. a name for an existing datastore in the same datacenter
1. a name for an existing resource pool in the cluster
1. a name for a Templates folder
1. a name for a VMs folder
1. a name for a disks folder

The folders mentioned above will be created during deployment if they do not exist at that time.