#!/bin/bash
GOOS=linux go build
docker build -t laziestperson1/reserve .
docker push laziestperson1/reserve
