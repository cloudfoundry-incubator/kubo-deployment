#!/bin/bash -e

[ -z "$DEBUG" ] || set -x

version="$(cat kubo-version/version)"
base_directory=git-kubo-deployment

pushd "${base_directory}"
  archive_name="kubo-deployment-${version}.tgz"
  git archive -o "${archive_name}" --prefix kubo-deployment/ HEAD

  echo "Generated ${archive_name}:"
  tar tvf "${archive_name}"
popd

mv "${base_directory}/${archive_name}" "tarballs/"
