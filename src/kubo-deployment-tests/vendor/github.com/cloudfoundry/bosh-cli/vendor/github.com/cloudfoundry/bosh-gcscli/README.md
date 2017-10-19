## GCS CLI

A CLI for uploading, fetching and deleting content to/from
the [GCS blobstore](https://cloud.google.com/storage/). This is **not**
an official Google Product.

## Installation

```
go get github.com/cloudfoundry/bosh-gcscli
```

## Usage

Given a JSON config file (`config.json`)...

``` json
{
  "bucket_name":         "name of GCS bucket (required)",

  "credentials_source":  "flag for credentials
                          (optional, defaults to Application Default Credentials)
                          (can be "static" for json_key),
                          (can be "none" for explicitly no credentials)"
  "storage_class":       "storage class for objects
                          (optional, defaults to bucket settings)",
  "json_key":            "JSON Service Account File
                          (optional, required for static credentials)",
  "encryption_key":      "Base64 encoded 32 byte Customer-Supplied
                          encryption key used to encrypt objects 
                          (optional)"
}
```


Empty `credentials_source` implies attempting to use Application Default
Credentials. `none` as `credentials_source` specifies no read-only scope
with explicitly no credentials. `static` as `credentials_source` specifies to
use the [Service Account File](https://developers.google.com/identity/protocols/OAuth2ServiceAccount) included
in `json_key`.

Empty `storage_class` implies using the default for the bucket.

``` bash
# Usage
bosh-gcscli --help

# Command: "put"
# Upload a blob to the GCS blobstore.
bosh-gcscli -c config.json put <path/to/file> <remote-blob>

# Command: "get"
# Fetch a blob from the GCS blobstore.
# Destination file will be overwritten if exists.
bosh-gcscli -c config.json get <remote-blob> <path/to/file>

# Command: "delete"
# Remove a blob from the GCS blobstore.
bosh-gcscli -c config.json delete <remote-blob>

# Command: "exists"
# Checks if blob exists in the GCS blobstore.
bosh-gcscli -c config.json exists <remote-blob>
```

Alternatively, this package's underlying client can be used to access GCS,
see the [godoc](https://godoc.org/github.com/cloudfoundry/bosh-gcscli)
for more information.

## Tooling

A Makefile is provided for ease of development. Targets are annotated
with descriptions.

gvt is used for vendoring. For full usage, see the [manual at godoc](https://godoc.org/github.com/FiloSottile/gvt).

Integration tests expect to be run from a host with [Application Default
Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
available which has permissions to create and delete buckets.
Application Default Credentials are present on any GCE instance and inherit
the permisions of the [service account](https://cloud.google.com/iam/docs/service-accounts)
assigned to the instance.

## License

This library is licensed under Apache 2.0. Full license text is
available in [LICENSE](LICENSE).