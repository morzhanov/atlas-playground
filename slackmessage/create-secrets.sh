#!/bin/bash

kubectl create secret generic -n argo webhook-token-secret \
  --from-literal=token=TOKEN
