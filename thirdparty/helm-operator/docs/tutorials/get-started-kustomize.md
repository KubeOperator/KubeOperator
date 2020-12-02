# How to bootstrap Helm operator using Kustomize

This guide shows you how to use Kustomize to bootstrap Helm Operator on a Kubernetes cluster.

## Prerequisites

You will need to have Kubernetes set up. For a quick local test,
you can use `minikube`, `kubeadm` or `kind`. Any other Kubernetes setup
will work as well though.

If working on e.g. GKE with RBAC enabled, you will need to add a cluster role binding:

```sh
kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
    --clusterrole=cluster-admin \
    --user="$(gcloud config get-value core/account)"
```

## Prepare Helm Operator installation 

Create a directory called `fluxcd` and add the `flux` namespace definition to it:

```sh
mkdir fluxcd

cat > fluxcd/namespace.yaml <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: flux
EOF
```

Create the `repositories.yaml` file and add the stable, flagger and podinfo Helm repositories to it:

```sh
cat > fluxcd/repositories.yaml <<EOF
apiVersion: v1
repositories:
- name: stable
  url: https://kubernetes-charts.storage.googleapis.com
  cache: /var/fluxd/helm/repository/cache/stable-index.yaml
- name: flagger
  url: https://flagger.app
  cache: /var/fluxd/helm/repository/cache/flagger-index.yaml
- name: podinfo
  url: https://stefanprodan.github.io/podinfo
  cache: /var/fluxd/helm/repository/cache/podinfo-index.yaml
EOF
```

Create a kustomization file and use the Helm operator deploy YAMLs as base:

```sh
cat > fluxcd/kustomization.yaml <<EOF
namespace: flux
resources:
  - namespace.yaml
bases:
 - github.com/fluxcd/helm-operator//deploy
secretGenerator:
  - name: helm-repositories
    files:
      - repositories.yaml
patchesStrategicMerge:
  - patch.yaml
EOF
```

> **Note:** If you want to install a specific Helm operator release,
> add the version number to the base URL:
> `github.com/fluxcd/helm-operator//deploy?ref=v1.0.0-rc2`


Create a patch file for Helm operator and mount the repositories secret:

```sh
cat > fluxcd/patch.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: flux-helm-operator
spec:
  template:
    spec:
      volumes:
       - name: repositories-yaml
         secret:
           secretName: helm-repositories
       - name: repositories-cache
         emptyDir: {}
      containers:
        - name: flux-helm-operator
          volumeMounts:
             - name: repositories-yaml
               mountPath: /var/fluxd/helm/repository
             - name: repositories-cache
               mountPath: /var/fluxd/helm/repository/cache
EOF
```

## Install Helm Operator with Kustomize

In the next step, deploy Flux to the cluster (you'll need kubectl **1.14** or newer):

```sh
kubectl apply -k fluxcd
```

Wait for Helm operator to start:

```sh
kubectl -n flux rollout status deployment/flux-helm-operator
```

## Use the `HelmRelease` custom resource

Install podinfo by referring to its Helm repository:

```sh
cat <<EOF | kubectl apply -f -
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: podinfo
  namespace: default
spec:
  releaseName: podinfo
  chart:
    repository: https://stefanprodan.github.io/podinfo
    version: 2.1.0
    name: podinfo
  values:
    replicaCount: 1
EOF
```

Verify that the Helm Operator has installed the release:

```sh
kubectl get hr

NAME       RELEASE    STATUS     MESSAGE                  AGE
podinfo    podinfo    DEPLOYED   helm install succeeded   1m
```

Delete the release with:

```sh
kubectl delete hr/podinfo
```

## Next steps

Try out [fluxcd/helm-operator-get-started](https://github.com/fluxcd/helm-operator-get-started)
to learn more about Helm Operator capabilities.
