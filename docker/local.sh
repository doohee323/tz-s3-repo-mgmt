#!/bin/bash

cd tz-s3-repo-mgmt/docker

sudo apt-get -y install docker.io
sudo usermod -G docker vagrant
sudo chown -Rf vagrant:vagrant /var/run/docker.sock
#service docker status

## for local env ###
docker-compose -f docker-compose.yml build --no-cache
docker image ls
docker-compose -f docker-compose.yml up --force-recreate
docker ps
#docker container stop $(docker container ls -a -q) && docker system prune -a -f --volumes

exit 0
