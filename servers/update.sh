#!/bin/bash
export RESERVE=reserve:80
export REDIS=redisServer:6379

docker rm -f gateway
docker rm -f redisServer
docker rm -f rabbit
docker rm -f reserve
docker network rm assignmentNetwork

docker network create assignmentNetwork
docker run -d --name redisServer --network assignmentNetwork redis
docker run -d --name rabbit --network assignmentNetwork rabbitmq:3

docker pull laziestperson1/reserve
docker pull laziestperson1/gateway


docker run -d \
    --network assignmentNetwork \
    -e ADDR=$RESERVE \
    --restart unless-stopped \
    --name reserve laziestperson1/reserve \
    ping rabbit


docker run -d \
--network assignmentNetwork \
-p 443:443 \
-e ADDR=:443 \
-e RESERVE=$RESERVE \
-e REDIS=$REDIS \
-e TLSCERT=/etc/letsencrypt/live/api.html-summary.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.html-summary.me/privkey.pem \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--restart unless-stopped \
--name gateway laziestperson1/gateway \
ping redisServer

