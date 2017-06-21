## BOSH Lite on VirtualBox

1. Check that your machine has at least 8GB RAM, and 100GB free disk space. Smaller configurations may work.

1. Install [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

    Known working version:

    ```
    $ VBoxManage --version
    5.1...
    ```

    Note: If you encounter problems with VirtualBox networking try installing [Oracle VM VirtualBox Extension Pack](https://www.virtualbox.org/wiki/Downloads) as suggested by [Issue 202](https://github.com/cloudfoundry/bosh-lite/issues/202). Alternatively make sure you are on VirtualBox 5.1+ since previous versions had a [network connectivity bug](https://github.com/concourse/concourse-lite/issues/9).

1. Install Director

    ```
    $ git clone https://github.com/cloudfoundry/bosh-deployment ~/workspace/bosh-deployment

    $ mkdir -p ~/deployments/vbox

    $ cd ~/deployments/vbox
    ```

    Below command will try automatically create/enable Host-only network 192.168.50.0/24 ([details](https://github.com/cppforlife/bosh-virtualbox-cpi-release/blob/master/docs/networks-host-only.md)) and NAT network 'NatNetwork' with DHCP enabled ([details](https://github.com/cppforlife/bosh-virtualbox-cpi-release/blob/master/docs/networks-nat-network.md)).

    ```
    $ bosh create-env ~/workspace/bosh-deployment/bosh.yml \
      --state ./state.json \
      -o ~/workspace/bosh-deployment/virtualbox/cpi.yml \
      -o ~/workspace/bosh-deployment/virtualbox/outbound-network.yml \
      -o ~/workspace/bosh-deployment/bosh-lite.yml \
      -o ~/workspace/bosh-deployment/bosh-lite-runc.yml \
      -o ~/workspace/bosh-deployment/jumpbox-user.yml \
      --vars-store ./creds.yml \
      -v director_name="Bosh Lite Director" \
      -v internal_ip=192.168.50.6 \
      -v internal_gw=192.168.50.1 \
      -v internal_cidr=192.168.50.0/24 \
      -v outbound_network_name=NatNetwork
    ```

1. Alias and log into the Director

    ```
    $ bosh -e 192.168.50.6 --ca-cert <(bosh int ./creds.yml --path /director_ssl/ca) alias-env vbox
    $ export BOSH_CLIENT=admin
    $ export BOSH_CLIENT_SECRET=`bosh int ./creds.yml --path /admin_password`
    ```

1. Upload BOSH Lite stemcell

    ```
    $ bosh -e vbox upload-stemcell https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent?v=3421.9 \
      --sha1 1396d7877204e630b9e77ae680f492d26607461d
    ```

1. Update cloud config

    ```
    $ bosh -e vbox update-cloud-config ~/workspace/bosh-deployment/warden/cloud-config.yml
    ```

1. Deploy example deployment

    ```
    $ bosh -e vbox -d zookeeper deploy <(wget -O- https://raw.githubusercontent.com/cppforlife/zookeeper-release/master/manifests/zookeeper.yml)
    ```

1. Optionally, set up a local route for `bosh ssh` command

    ```
    $ sudo route add -net 10.244.0.0/16    192.168.50.6 # Mac OS X
    $ sudo route add -net 10.244.0.0/16 gw 192.168.50.6 # Linux
    $ route add           10.244.0.0/16    192.168.50.6 # Windows
    ```

1. In case you need to SSH into the Director VM, see [Jumpbox user](jumpbox-user.md).

1. In case VirtualBox VM shuts down or reboots

    To bring back the Director on the VM you will have to re-run `create-env` command from above. (You will have to remove `current_manifest_sha` line from `state.json` to force a redeploy.) The containers will be lost after a VM restart, but you can restore your deployment with `bosh cck` command. Alternatively *Pause* the VM from the VirtualBox UI before shutting down VirtualBox host, or making your computer sleep.
