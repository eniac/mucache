#!/bin/bash -ex

export MUCACHE_TOP=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

## Make sure that the containers exist and have been built first
export docker_io_username=${1?docker.io username not given}
export application_namespace=${2?application name not given, e.g., social}
export cm_enabled=${3:-true}
export shard=${4:-1}

declare -A all_services
all_services["social"]="post_storage home_timeline user_timeline social_graph compose_post"
all_services["twoservices"]="caller callee"

services=(${all_services[$application_namespace]})

## Services
for idx in "${!services[@]}"; do
  app_name=${services[$idx]}
  app_name_no_underscores=${app_name//_/}
  node_idx=$((idx + 1))
  for shard_idx in $(seq 1 "$shard"); do
    NODE_IDX="${node_idx}" \
      SHARD_IDX="$shard_idx" \
      SHARD_COUNT="$shard" \
      CM_ENABLED="$cm_enabled" \
      APP_NAMESPACE="$application_namespace" \
      APP_NAME="$app_name""$shard_idx" \
      APP_RAW_NAME="$app_name" \
      APP_RAW_NAME_NO_UNDERSCORES="$app_name_no_underscores" \
      APP_NAME_NO_UNDERSCORES="$app_name_no_underscores""$shard_idx" \
      envsubst <"${MUCACHE_TOP}/deploy/shard/app.yaml" | kubectl apply -f -
  done
done

if [ "$cm_enabled" = "true" ]; then
  ## Cache Manager
  for idx in "${!services[@]}"; do
    for shard_idx in $(seq 1 "$shard"); do
      cm_adds="/app/experiments/k8s_cm/$application_namespace.txt"
      node_idx=$((idx + 1))
      NODE_IDX="${node_idx}" \
        SHARD_IDX="$shard_idx" \
        SHARD_COUNT="$shard" \
        CM_ADDS=$cm_adds \
        envsubst <"${MUCACHE_TOP}/deploy/shard/cm/cm.yaml" | kubectl apply -f -
    done
  done
fi
