# Troubleshooting

## Problem

The error below occurs when deploying KuBOSH on GCP

`CmdError{"type":"Bosh::Clouds::CloudError","message":"Creating stemcell: Creating Google Storage Bucket: Post https://www.googleapis.com/storage/v1/b?alt=json\u0026project=cf-pcf-kubo: Get http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token: dial tcp 169.254.169.254:80: i/o timeout","ok_to_retry":false}`

### Solution

Make sure that you have logged in with the newly created service account as described in 
[GCP Setup](https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/tree/master/docs/bosh#setup).

## Problem

When deploying KuBOSH via an `sshuttle` connection, the following error might occur:

```
Command 'deploy' failed:
  Deploying:
    Creating instance 'bosh/0':
      Waiting until instance is ready:
        Sending ping to the agent:
          Performing request to agent endpoint 'https://mbus:294a691d057ede1af4f696aab36c4bc5@10.0.0.4:6868/agent':
            Performing POST request:
              Post https://mbus:294a691d057ede1af4f696aab36c4bc5@10.0.0.4:6868/agent: dial tcp 10.0.0.4:6868: i/o timeout
```

where `10.0.0.4` is the KuBOSH IP address.

### Solution

Restart the `sshuttle` connection.

## Problem

A KuBOSH or kubo deployment may fail with the following error:

```
The provided username and password combination are incorrect. Please validate your input and retry your request.
```

### Solution

This is typically caused by version mismatch between the Credhub CLI and backend. The versions must match exactly.
Please make sure that the installed Credhub CLI matches version in the 
[requirements section](../README.md#required-software):

```
$ credhub --version
CLI Version: 0.4.0
Server Version: 0.4.0
```

## Problem

Various connectivity issues during deployment

### Solution

Please make sure that all the host names used in the configuration are resolving to the correct IP addresses.

## Problem

Strange error when running bosh-cli 

```
bosh-cli -e p-kubo deployments
Using environment 'p-kubo.p.kubo.cf-app.com' as '?'

Finding deployments:
  Performing request GET 'https://p-kubo.p.kubo.cf-app.com:25555/deployments':
    Performing GET request:
      Refreshing token: UAA responded with non-successful status code '401' response '{"error":"invalid_token","error_description":"The token expired, was revoked, or the token ID is incorrect: acc3bc0eeb2b4323995f6d5873d9f52e-r"}'

Exit code 1
```

### Solution

Please make sure that you are logged in. The password can be found in `<BOSH_ENV>/creds.yml`, in the `admin_password` field.
