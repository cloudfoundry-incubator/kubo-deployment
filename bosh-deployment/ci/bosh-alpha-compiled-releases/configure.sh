#!/usr/bin/env bash

set -eu

fly -t production set-pipeline \
 -p bosh-alpha-compiled-releases \
 -c ./pipeline.yml \
 -l <(lpass show --note "concourse:production pipeline:bosh-alpha-compiled-releases")
