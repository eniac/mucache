#!/bin/bash -ex
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

if [ $# -eq 0 ]; then
  echo "Please provide the number of worker nodes"
  exit 1
fi

nworkers=$1

setup_path="$HOME"/mucache/scripts/setup

# Launch Redis on each worker node
for i in $(seq 1 "$nworkers"); do
  sed "s/node-1/node-$i/g" "$setup_path"/redis.yaml >"$setup_path"/redis"$i".yaml
#  helm install redis"$i" bitnami/redis -f "$setup_path"/redis"$i".yaml
  helm install redis"$i" oci://registry-1.docker.io/bitnamicharts/redis -f "$setup_path"/redis"$i".yaml
  rm "$setup_path"/redis"$i".yaml
done

# Launch Memcached on each worker node
#for i in $(seq 1 "$nworkers"); do
#  sed "s/node-1/node-$i/g" "$setup_path"/memcached.yaml >"$setup_path"/memcached"$i".yaml
#  helm install memcached"$i" bitnami/memcached -f "$setup_path"/memcached"$i".yaml
#  rm "$setup_path"/memcached"$i".yaml
#done
for i in $(seq 1 "$nworkers"); do
  NODE_IDX="$i" \
    MEM=0 \
    envsubst <"$setup_path"/cache.yaml | helm install cache"$i" bitnami/redis -f -
done

# Initialize dapr
dapr init --kubernetes --wait --enable-mtls=false --runtime-version 1.10.4

# Launch Redis Component
for i in $(seq 1 "$nworkers"); do
  sed "s/redis1/redis$i/g" "$setup_path"/dapr_redis.yaml >"$setup_path"/dapr_redis"$i".yaml
  kubectl apply -f "$setup_path"/dapr_redis"$i".yaml
  rm "$setup_path"/dapr_redis"$i".yaml
done

kubectl apply -f "$setup_path"/config.yaml

# Metrics
#kubectl create namespace dapr-monitoring
#helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
#helm repo update
#helm install dapr-prom prometheus-community/prometheus -n dapr-monitoring --set alertmanager.persistentVolume.enable=false --set pushgateway.persistentVolume.enabled=false --set server.persistentVolume.enabled=false
#
#helm repo add grafana https://grafana.github.io/helm-charts
#helm repo update
#helm install grafana grafana/grafana -n dapr-monitoring --set persistence.enabled=false
# https://docs.dapr.io/operations/monitoring/metrics/grafana/
