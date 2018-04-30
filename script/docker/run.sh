#!/bin/bash


url=http://10.10.200.97:28080/cpuwork.php/?cpu=10
options="--v=3 --target=$url --logtostderr --rps=2"

#docker run -d beekman9527/webclient:v1 $options 
./_output/webClient $options

sleep 1
docker ps
