#!/bin/bash

export MUCACHE_TOP=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

export docker_io_username=${1?docker.io username not given}

tag="${docker_io_username}/mucache"

docker build \
  -f "${MUCACHE_TOP}/deploy/Dockerfile" \
  -t "${tag}" .
docker push "${tag}"
