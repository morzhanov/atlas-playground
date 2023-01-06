# Setting up the GreenOps PoC

- Follow the steps on the [getting started](https://greenops.io/go-docs/getting-started/) page to set up GreenOps.
- Create a quick secret in the Argo CD namespace: ```kubectl create secret generic greenops-sso --from-literal=client-id=Z3JlZW5vcHMtc3Nv --from-literal=client-secret=<SSO_KEY> -n argocd``` and in the GreenOps namespace: ```kubectl create secret generic greenops-sso --from-literal=client-id=Z3JlZW5vcHMtc3Nv --from-literal=client-secret=<SSO_KEY> -n greenops```
- Create a secret containing the admin login for Argo CD for **each** cluster in the `argo` namespace. The keys should be `username` and `password`.
  ```
  ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo)
  kubectl create secret generic argocd-user --from-literal=username=admin --from-literal=password=$ARGOCD_PASSWORD -n argo
  ```
- *Many of these steps will be automated in the future*

# Setting up the Emailtree pipeline
- Create `github-pat` token in the `argo` namespace for pulling info from Git. Should have the key `token`. Example command: ```kubectl create secret generic github-pat --from-literal=token=<GIT_TOKEN> -n argo```
- Install  templates from the template section using the filenames.
- Update the variables in `emailtree/poc/pipeline/params.yaml`.
    - `repository` should point to your own repo.
    - `filepath` should point to wherever the app/cluster information is stored.
- Set up 4 Argo applications. For this example, they can all be in one environment. They should be named `rollout-env1`, `rollout-env2`, `rollout-env3`, `rollout-env4`.
- Create namespaces for the Argo CD apps:
  ```
  kubectl create ns env1
  kubectl create ns env2
  kubectl create ns env3
  kubectl create ns env4
  ```
- Run the pipeline!

**Note**: This pipeline can easily be run with different names, across different environments, and with a different repo structure. This was written simply so the example will work as quickly as possible.

# What does the pipeline do?
- The pipeline dynamically fetches the list of all applications (across clusters) that needs to be deployed or synced.
- The deployments are made in parallel (with a limit of 3). Only 3 deployments can run in parallel at any given time. If a deployment fails, the workflow will not continue with other deployments.
- Deployment process:
  - Enable maintenance mode using a Job
  - Deploy new application using a blue-green Argo Rollout
  - Disable maintenance mode using a Job
