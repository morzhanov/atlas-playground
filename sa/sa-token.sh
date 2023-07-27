# 1. Create the token
TOKEN=$(openssl rand -base64 32)

# 2. Define the secret in a YAML file
echo "
apiVersion: v1
kind: Secret
metadata:
  name: workflow.service-account-token
  namespace: argo
  annotations:
    kubernetes.io/service-account.name: workflow
    kubernetes.io/service-account.uid: $(kubectl get serviceaccounts workflow -n argo -o jsonpath='{.metadata.uid}')
type: kubernetes.io/service-account-token
data:
  token: $(echo -n $TOKEN | base64 | tr -d '\n')
  ca.crt: $(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
  namespace: $(kubectl config view --raw --minify --flatten -o jsonpath='{.contexts[].context.namespace}' | base64 | tr -d '\n')
" > secret.yml

# 3. Apply the secret
kubectl apply -f secret.yml
