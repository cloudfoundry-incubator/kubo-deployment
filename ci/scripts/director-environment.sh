#!/bin/sh

creds_path="$1/creds.yml"
export BOSH_CLIENT="bosh-admin"
export BOSH_CLIENT_SECRET="$(bosh-cli int "$creds_path" --path /admin_password)"
