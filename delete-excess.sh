#!/bin/bash

for ns in $(seq 1 35); do
  namespace="ns-$ns"
  for i in 4; do
    echo deleteing deployment fortio-client-$i and fortio-server-$i in $namespace
    kubectl delete deployment fortio-client-$i -n $namespace 
    kubectl delete deployment fortio-server-$i -n $namespace 
    kubectl delete svc fortio-server-service-$i -n $namespace 
  done
done
