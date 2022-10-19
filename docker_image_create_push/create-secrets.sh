#!/bin/bash

kubectl create secret docker-registry dockercred -n argo \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=USERNAME \
  --docker-password=PASSOWRD \
  --docker-email=EMAIL

kubectl create secret generic -n argo git-cred-secret \
  --from-literal=username=USERNAME \
  --from-literal=token=TOKEN
