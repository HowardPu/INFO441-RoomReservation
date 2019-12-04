#!/bin/bash
npm run build
docker build -t laziestperson1/reservation-client .
docker push laziestperson1/reservation-client
