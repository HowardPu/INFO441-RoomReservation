#!/bin/bash
export MESSAGE=messaging:80
export MESSAGE1=messaging1:80
export SUMMARY=summary:80
export REDIS=redisServer:6379

docker rm -f gateway
docker rm -f redisServer
docker rm -f summary
docker rm -f messaging
docker rm -f messaging1
docker rm -f rabbit
docker network rm assignmentNetwork

docker network create assignmentNetwork
docker run -d --name redisServer --network assignmentNetwork redis
docker run -d --name rabbit --network assignmentNetwork rabbitmq:3


docker pull zkzhiqilin/messaging

docker run -d \
--network assignmentNetwork \
-e ADDR=$MESSAGE \
--restart unless-stopped \
--name messaging zkzhiqilin/messaging


docker run -d \
--network assignmentNetwork \
-e ADDR=$MESSAGE1 \
--restart unless-stopped \
--name messaging1 zkzhiqilin/messaging


docker pull zkzhiqilin/gateway
docker run -d \
--network assignmentNetwork \
-p 443:443 \
-e ADDR=:443 \
-e MESSAGE=$MESSAGE \
-e MESSAGE1=$MESSAGE1 \
-e SUMMARY=$SUMMARY \
-e REDIS=$REDIS \
-e TLSCERT=/etc/letsencrypt/live/api.awesome-summary.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.awesome-summary.me/privkey.pem \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--restart unless-stopped \
--name gateway zkzhiqilin/gateway \
ping redisServer

