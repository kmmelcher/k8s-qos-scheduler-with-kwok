#!/bin/bash

# Installing docker
sudo apt-get update
sudo apt-get install docker.io

DOCKER_VERSION=`docker --version`

echo "Intalled Docker version: "$DOCKER_VERSION

# starting and enabling Docker
sudo systemctl enable docker