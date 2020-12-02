# Upgrading from Helm operator beta (`>=0.5.0`) to stable (`>=1.0.0`)

Due to the Flux CD project joining the CNCF Sandbox and the API
becoming stable, the Helm operator has undergone changes that
necessitate some changes to your `HelmRelease` resources.

The central difference is that the Helm operator now works with
resources of the kind `HelmRelease` in the API version
`helm.fluxcd.io/v1`, the format of the resource is backwards
compatible.

Here are some things to know:

- The new operator will ignore the old custom resources (and the old
  operator will ignore the new resources).
- Deleting a resource while the corresponding operator is running
  will result in the Helm release also being deleted
- Deleting a `CustomResourceDefinition` will also delete all
  custom resources of that kind.
- If both operators are running and both new and old custom resources
  defining a release, the operators will fight over the release.
  
The safest way to upgrade is to avoid deletions and fights by stopping
the old operator. Replacing it with the new one (e.g., by changing the
deployment, or re-releasing the Flux chart with the new version) will
have that effect.

Once the old operator is not running, it is safe to deploy the new
operator, and start replacing the old resources with new
resources. You can keep the old resources around during this process,
since the new operator will ignore them.

## Updating custom resources

> **Note:** once the new CRD is applied it is no longer possible
> to list your old and new `HelmRelease` resources with just `kubectl
> get <hr|helmrelease>`, due to them sharing the same names. It is
> however still possible to list them by their full name.
>
> ```bash
> # Old `HelmRelease` resources
> $ kubectl get helmreleases.flux.weave.works
> # New `HelmRelease` resources
> $ kubectl get helmreleases.helm.fluxcd.io
> ```

The only difference between the old resource format and the new is
the changed API version.

Changing an old resource to a new resource is thus as simple as
changing the `apiVersion` field to `helm.fluxcd.io/v1`.

As a full example, this is an old resource:

```yaml
---
apiVersion: flux.weave.works/v1beta1
kind: HelmRelease
metadata:
  name: foobar
  namespace: foo-ns
spec:
 chart:
    git: git@example.com:user/repo
    path: charts/foobar
  values:
    image:
      repository: foobar
      tag: v1
```

The new custom resource would be:

```yaml
---
apiVersion: helm.fluxcd.io/v1       # <- change API version
kind: HelmRelease
metadata:
  name: foobar
  namespace: foo-ns
spec:
 chart:
    git: git@example.com:user/repo
    path: charts/foobar
  values:
    image:
      repository: foobar
      tag: v1
```

## Deleting the old resources

Once you have migrated all your `HelmRelease` resources to the new API
version and domain. You can remove all of the old resources by removing
the old Custom Resource Definition.

```sh
$ kubectl delete crd helmreleases.flux.weave.works
```
