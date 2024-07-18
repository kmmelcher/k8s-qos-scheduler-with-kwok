#!/bin/bash

# Add Kubernetes Signing Key.
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add

# Kubernetes is not included in the default repositories. Adding the Kubernetes reposiroty.
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
sudo apt-get update

# Installing the Kubernetes tools.
sudo apt-get install kubeadm kubelet kubectl
sudo apt-mark hold kubeadm kubelet kubectl
kubeadm version