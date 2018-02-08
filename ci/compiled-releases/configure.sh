#!/usr/bin/env bash

set -eu

fly -t production set-pipeline -n \
 -p compiled-releases-3468 \
 -c ./pipeline-3468.yml \
 -l <(lpass show --note "concourse:production pipeline:compiled-releases")

fly -t production check-resource -r compiled-releases-3468/bosh-release -f version:263.4.0
fly -t production check-resource -r compiled-releases-3468/uaa-release -f version:52.2
# fly -t production check-resource -r compiled-releases-3468/credhub-release -f version:1.6.0
# fly -t production check-resource -r compiled-releases-3468/warden-cpi -f version:37
# fly -t production check-resource -r compiled-releases-3468/garden-linux -f version:0.342.0
# fly -t production check-resource -r compiled-releases-3468/garden-runc -f version:1.9.4
# fly -t production check-resource -r compiled-releases-3468/grootfs -f version:0.24.0
fly -t production check-resource -r compiled-releases-3468/ubuntu-trusty-stemcell -f version:3468

fly -t production set-pipeline -n \
 -p compiled-releases-3445 \
 -c ./pipeline-3445.yml \
 -l <(lpass show --note "concourse:production pipeline:compiled-releases")

fly -t production check-resource -r compiled-releases-3445/bosh-release -f version:263
fly -t production check-resource -r compiled-releases-3445/uaa-release -f version:45.4
# fly -t production check-resource -r compiled-releases-3445/credhub-release -f version:1.3.4
fly -t production check-resource -r compiled-releases-3445/ubuntu-trusty-stemcell -f version:3445

fly -t production set-pipeline -n \
 -p compiled-releases-3421 \
 -c ./pipeline-3421.yml \
 -l <(lpass show --note "concourse:production pipeline:compiled-releases")

fly -t production check-resource -r compiled-releases-3421/bosh-release -f version:262.4
fly -t production check-resource -r compiled-releases-3421/uaa-release -f version:41.1
# fly -t production check-resource -r compiled-releases-3421/credhub-release -f version:1.0.8
fly -t production check-resource -r compiled-releases-3421/ubuntu-trusty-stemcell -f version:3421

fly -t production set-pipeline -n \
 -p compiled-releases-3363 \
 -c ./pipeline-3363.yml \
 -l <(lpass show --note "concourse:production pipeline:compiled-releases")

fly -t production check-resource -r compiled-releases-3363/bosh-release -f version:260.8
fly -t production check-resource -r compiled-releases-3363/uaa-release -f version:24.12
fly -t production check-resource -r compiled-releases-3363/ubuntu-trusty-stemcell -f version:3363
