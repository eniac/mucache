#!/bin/bash -ex

# disable swap
sudo swapoff -a

# use bash as default shell
sudo usermod -s /bin/bash "$USER"

# oha
echo "deb [signed-by=/usr/share/keyrings/azlux-archive-keyring.gpg] http://packages.azlux.fr/debian/ stable main" | sudo tee /etc/apt/sources.list.d/azlux.list
sudo wget -O /usr/share/keyrings/azlux-archive-keyring.gpg https://azlux.fr/repo.gpg
# python3, we need 3.7+
sudo add-apt-repository ppa:deadsnakes/ppa

sudo apt update
sudo apt install -y libssl-dev zlib1g-dev
sudo apt install -y mosh htop redis-tools oha python3.10 python3-pip

pip3 install requests progress tqdm

# rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >>~/.bashrc

# dapr
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash

# wrk
git clone https://github.com/giltene/wrk2.git
(
  cd wrk2 || exit
  make -j
  sudo apt install -y luarocks
  sudo luarocks install https://raw.githubusercontent.com/tiye/json-lua/main/json-lua-0.1-4.rockspec
  sudo luarocks install luasocket
  sudo luarocks install uuid
)

# helm v3
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
rm get_helm.sh

# Redis helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# k3s
HOSTNAME=$(hostname | cut -d. -f1)
curl -sfL https://get.k3s.io | sh -s - --token 111 --write-kubeconfig-mode 644 --node-name "$HOSTNAME" --flannel-backend=host-gw --disable traefik,servicelb --disable-network-policy

echo "export PYTHONPATH=$HOME/mucache/" >>"$HOME"/.bashrc
echo "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml" >>"$HOME"/.bashrc

echo "Remember to run" '`source ~/.bashrc`' "to reload the path!"
