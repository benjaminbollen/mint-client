#! /bin/sh

cd $GOPATH/src/github.com/eris-ltd/mint-client
docker build -t mcta_client -f ./DOCKER/all_in_one/Dockerfile . 
docker run -t --rm mcta_client
