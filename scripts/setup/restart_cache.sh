#!/bin/bash -ex

if [ $# -eq 0 ]; then
  echo "Please provide the number of worker nodes"
  exit 1
fi

nworkers=$1
mem=$2

setup_path="$HOME"/mucache/scripts/setup

for i in $(seq 1 "$nworkers"); do
  helm uninstall cache"$i" || true
done

#for i in $(seq 1 "$nworkers"); do
#  NODE_IDX="$i" \
#    MEM="$mem" \
#    envsubst <"$setup_path"/cache.yaml | helm install cache"$i" bitnami/redis -f -
#done
for i in $(seq 1 "$nworkers"); do
  NODE_IDX="$i" \
    MEM="$mem" \
    envsubst <"$setup_path"/cache.yaml | helm install cache"$i" oci://registry-1.docker.io/bitnamicharts/redis -f -
done
