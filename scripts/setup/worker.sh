#!/bin/bash -ex

# disable swap
sudo swapoff -a

HOSTNAME=$(hostname | cut -d. -f1)

curl -sfL https://get.k3s.io | K3S_URL=https://10.10.1.1:6443 sh -s - agent --token 111 --node-name "$HOSTNAME"
