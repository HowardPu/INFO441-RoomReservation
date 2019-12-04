#!/bin/bash
docker rm -f reservation-client
docker pull laziestperson1/reservation-client
docker run -d -p 443:443 -p 80:80 -e ADDR=:443 -v /etc/letsencrypt:/etc/letsencrypt:ro --name reservation-client laziestperson1/reservation-client  