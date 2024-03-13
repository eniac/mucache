#!/bin/bash -ex

if [ $# -eq 0 ]; then
  echo "Please provide the number of worker nodes"
  exit 1
fi

nworkers=$1
mem=$2

setup_path="$HOME"/mucache/scripts/setup

for i in $(seq 1 "$nworkers"); do
  helm uninstall memcached"$i"
done

if [ -z "$mem" ]; then
  for i in $(seq 1 "$nworkers"); do
    sed "s/node-1/node-$i/g" "$setup_path"/memcached.yaml >"$setup_path"/memcached"$i".yaml
    helm install memcached"$i" bitnami/memcached -f "$setup_path"/memcached"$i".yaml
    rm "$setup_path"/memcached"$i".yaml
  done
else
  # use a different memcached.yaml file for now
  # because we don't know if setting the memory manually change other setting
  for i in $(seq 1 "$nworkers"); do
    sed "s/node-1/node-$i/g" "$setup_path"/memcached_size.yaml >"$setup_path"/memcached"$i".yaml
    MEM="$mem" \
      envsubst <"$setup_path"/memcached"$i".yaml | helm install memcached"$i" bitnami/memcached -f -
    rm "$setup_path"/memcached"$i".yaml
  done
fi
