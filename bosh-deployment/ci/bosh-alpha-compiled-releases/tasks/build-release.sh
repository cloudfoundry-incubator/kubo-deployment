#!/usr/bin/env bash

set -eu

bosh -n create-release --dir bosh-src --tarball source-release/release.tgz
