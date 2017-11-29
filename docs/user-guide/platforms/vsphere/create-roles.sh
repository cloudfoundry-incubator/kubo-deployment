#!/bin/bash

set -ex

govc role.create manage-k8s-node-vms \
    Resource.AssignVMToPool \
    System.Anonymous \
    System.Read \
    System.View \
    VirtualMachine.Config.AddExistingDisk \
    VirtualMachine.Config.AddNewDisk \
    VirtualMachine.Config.AddRemoveDevice \
    VirtualMachine.Config.RemoveDisk \
    VirtualMachine.Inventory.Create \
    VirtualMachine.Inventory.Delete	

govc role.create manage-k8s-volumes \
    Datastore.AllocateSpace \
    Datastore.FileManagement \
    System.Anonymous \
    System.Read \
    System.View

govc role.create k8s-system-read-and-spbm-profile-view \
    StorageProfile.View \
    System.Anonymous \
    System.Read \
    System.View
