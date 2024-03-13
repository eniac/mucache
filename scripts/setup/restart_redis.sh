#!/bin/bash -ex

if [ $# -eq 0 ]; then
  echo "Please provide the number of worker nodes"
  exit 1
fi

nworkers=$1

setup_path="$HOME"/mucache/scripts/setup

for i in $(seq 1 "$nworkers"); do
  helm uninstall redis"$i" || true
done

for i in $(seq 1 "$nworkers"); do
  sed "s/redis1/redis$i/g" "$setup_path"/dapr_redis.yaml >"$setup_path"/dapr_redis"$i".yaml
  kubectl delete -f "$setup_path"/dapr_redis"$i".yaml || true
  rm "$setup_path"/dapr_redis"$i".yaml
done

for i in $(seq 1 "$nworkers"); do
  sed "s/node-1/node-$i/g" "$setup_path"/redis.yaml >"$setup_path"/redis"$i".yaml
  helm install redis"$i" bitnami/redis -f "$setup_path"/redis"$i".yaml
  rm "$setup_path"/redis"$i".yaml
done

for i in $(seq 1 "$nworkers"); do
  sed "s/redis1/redis$i/g" "$setup_path"/dapr_redis.yaml >"$setup_path"/dapr_redis"$i".yaml
  kubectl apply -f "$setup_path"/dapr_redis"$i".yaml
  rm "$setup_path"/dapr_redis"$i".yaml
done
