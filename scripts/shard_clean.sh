#!/bin/bash -x

export MUCACHE_TOP=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

export application_namespace=${1?application name not given, e.g., social}
export shard=${2:-1}

declare -A all_services
all_services["social"]="post_storage home_timeline user_timeline social_graph compose_post"
all_services["twoservices"]="caller callee"

## Services
for app_name in ${all_services[$application_namespace]}; do
  app_name_no_underscores=${app_name//_/}
  for shard_idx in $(seq 1 "$shard"); do
    APP_NAME_NO_UNDERSCORES="$app_name_no_underscores""$shard_idx" \
      envsubst <"${MUCACHE_TOP}/deploy/shard/app.yaml" | kubectl delete -f -
  done
done

services=(${all_services[$application_namespace]})

## Cache Manager
for idx in "${!services[@]}"; do
  for shard_idx in $(seq 1 "$shard"); do
    NODE_IDX=$((idx + 1)) \
      SHARD_IDX="$shard_idx" \
      envsubst <"${MUCACHE_TOP}/deploy/shard/cm/cm.yaml" | kubectl delete -f -
  done
done
