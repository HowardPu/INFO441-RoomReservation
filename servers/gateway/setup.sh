#!/bin/bash
GOOS=linux go build
docker build -t zkzhiqilin/gateway .
docker push zkzhiqilin/gateway
