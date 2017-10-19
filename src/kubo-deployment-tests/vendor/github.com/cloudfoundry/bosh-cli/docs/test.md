## Unit Tests

Each package in the CLI has its own unit tests and there are integration tests in the `integration` package.

You can also run all tests with `bin/test`.

## Acceptance Tests

The acceptance tests are designed to exercise the main commands of the CLI (deployment, deploy, delete).

They are not designed to verify the compatibility of CPIs or testing BOSH releases.

The acceptance test related to compiled releases uses an already compiled release that was compiled against a stemcell 
with os/version : ubuntu-trusty/2776. If the stemcell os/version used for the tests changes, you will need to modify the 
file [acceptance/assets/sample-release-compiled.tgz](acceptance/assets/sample-release-compiled.tgz). This can be either 
modified manually by un-zipping it and changing the release manifest, or you can refer build the sample release from 
[acceptance/assets/sample-release](acceptance/assets/sample-release) folder, upload it to a bosh installation, and 
deploy it against a stemcell with the desired OS and Version. Then use bosh export to export a compiled release.

### Fly executing the acceptance tests

In theory you should be able to export the environment variables for the task,
but I've had trouble getting that to work.

A way to run the acceptance tests that seems to work:

We're also going to need all the inputs for the task. `bosh-cli` is easy,
that's going to be the bosh-cli source directory. To satisfy the
`bosh-warden-cpi-release` input, we'll need to download the warden cpi release
(probably from bosh.io) and name it `cpi-release.tgz`.

Now we can fly execute:

```
./fly -t bosh-init -k execute -p -c ci/tasks/test-acceptance.yml -i bosh-cli=<path-to-source-dir> -i bosh-warden-cpi-release=<path-to-dir-containing-cpi-release.tgz>
```

