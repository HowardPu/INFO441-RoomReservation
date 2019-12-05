#!/bin/bash
GOOS=linux go build
docker build -t laziestperson1/gateway .
docker push laziestperson1/gateway
