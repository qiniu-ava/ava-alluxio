#!/bin/bash

set -ex

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

DATA_DISK=/disk1

# ===============================================================
# install git
# ---------------------------------------------------------------
apt-get update
apt-get install -y git ca-certificates curl

# ===============================================================
# install ansible
# ---------------------------------------------------------------

sudo apt-get install -y software-properties-common
sudo apt-add-repository ppa:ansible/ansible
sudo apt update
sudo apt-get install -y ansible

# ===============================================================
# install docker
# ---------------------------------------------------------------

sudo cp ${DIR}/daemon.json /etc/docker/daemon.json
sudo apt-get install -y apt-transport-https
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce

# ===============================================================
# docker settings
# ---------------------------------------------------------------

sudo groupadd docker
sudo usermod -aG docker $USER
sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
sudo chmod g+rwx "/home/$USER/.docker" -R

sudo service docker stop
mkdir -p -m 777 $DATA_DISK/docker

tar -zcC /var/lib docker > $DATA_DISK/var_lib_docker-backup-$(date +%s).tar.gz
sudo mv /var/lib/docker $DATA_DISK/docker
ln -s $DATA_DISK/docker /var/lib/docker
sudo service docker start

sudo systemctl enable docker

# ===============================================================
# pull basic docker images
# ---------------------------------------------------------------
docker pull ubuntu:16.04
