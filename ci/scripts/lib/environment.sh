#!/bin/sh

export KUBO_DEPLOYMENT_DIR="$(cd "$(dirname "$0")/../.."; pwd)"
export KUBO_ENVIRONMENT_DIR="${KUBO_DEPLOYMENT_DIR}/ci/environments/$(cat kubo-lock/name)"
