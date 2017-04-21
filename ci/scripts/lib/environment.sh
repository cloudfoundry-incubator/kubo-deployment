#!/bin/sh

export KUBO_DEPLOYMENT_DIR="$(cd "$(dirname "$0")/../.."; pwd)"
export KUBO_ENVIRONMENT_DIR="${PWD}/environment"
mkdir -p "${KUBO_ENVIRONMENT_DIR}"
echo "gcp" > "${KUBO_ENVIRONMENT_DIR}"/iaas
