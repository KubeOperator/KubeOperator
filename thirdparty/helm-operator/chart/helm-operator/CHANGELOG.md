## 0.6.0 (2020-01-26)

### Improvements

 - Update Helm Operator to `1.0.0-rc8`
   [fluxcd/helm-operator#244](https://github.com/fluxcd/helm-operator/pull/244)
 - Allow pod annotations, labels and account annotations to be set
   [fluxcd/helm-operator#229](https://github.com/fluxcd/helm-operator/pull/229)

## 0.5.0 (2020-01-10)

### Improvements

 - Update Helm Operator to `1.0.0-rc7`
   [fluxcd/helm-operator#197](https://github.com/fluxcd/helm-operator/pull/197)
 - Add support for configuring cert files for repositories
   [fluxcd/helm-operator#183](https://github.com/fluxcd/helm-operator/pull/183)
 - Add support for configuring Helm v3 repositories
   [fluxcd/helm-operator#173](https://github.com/fluxcd/helm-operator/pull/173)
 - Add Prometheus Operator ServiceMonitor templates
   [fluxcd/helm-operator#139](https://github.com/fluxcd/helm-operator/pull/139)

## 0.4.0 (2019-12-23)

### Improvements

 - Add `helm.versions` option to chart values
   [fluxcd/helm-operator#159](https://github.com/fluxcd/helm-operator/pull/159)
 - Update Helm Operator to `1.0.0-rc5`
   [fluxcd/helm-operator#157](https://github.com/fluxcd/helm-operator/pull/157)
 - Add Service and ServiceMonitor templates
   [fluxcd/helm-operator#139](https://github.com/fluxcd/helm-operator/pull/139)
 - Add extraVolumes and extraVolumeMounts
   [fluxcd/helm-operator#125](https://github.com/fluxcd/helm-operator/pull/125)

## 0.3.0 (2019-11-22)

### Improvements

 - Update Helm Operator to `1.0.0-rc4`
   [fluxcd/helm-operator#114](https://github.com/fluxcd/helm-operator/pull/114)
 - Fix upgrade command install instructions in `README.md`
   [fluxcd/helm-operator#92](https://github.com/fluxcd/helm-operator/pull/92)
 - Add `git.defaultRef` option for configuring an alternative Git default ref
   [fluxcd/helm-operator#83](https://github.com/fluxcd/helm-operator/pull/83)
 - Allow for deploying Tiller as a sidecar by setting `tillerSidecar.enabled`
   [fluxcd/helm-operator#79](https://github.com/fluxcd/helm-operator/pull/79)

## 0.2.1 (2019-10-18)

### Improvements

 - Update Helm Operator to `1.0.0-rc3`
   [fluxcd/helm-operator#74](https://github.com/fluxcd/helm-operator/pull/74)

## 0.2.0 (2019-10-07)

### Improvements

 - Update Helm Operator to `1.0.0-rc2`
   [fluxcd/helm-operator#59](https://github.com/fluxcd/helm-operator/pull/59)
 - Expand the list of public Helm repositories in the default config
   [fluxcd/helm-operator#53](https://github.com/fluxcd/helm-operator/pull/53)
 - Add `statusUpdateInterval` option for configuring the interval at which the operator consults Tiller for the status of a release
   [fluxcd/helm-operator#44](https://github.com/fluxcd/helm-operator/pull/44)

## 0.1.1 (2019-09-15)

### Improvements

 - Restart operator on helm repositories changes
   [fluxcd/helm-operator#30](https://github.com/fluxcd/helm-operator/pull/30)
 - Add liveness and readiness probes
   [fluxcd/helm-operator#30](https://github.com/fluxcd/helm-operator/pull/30)
 - Add `HelmRelease` example to chart notes
   [fluxcd/helm-operator#30](https://github.com/fluxcd/helm-operator/pull/30)

### Bug fixes

 - Fix SSH key mapping
   [fluxcd/helm-operator#30](https://github.com/fluxcd/helm-operator/pull/30)

## 0.1.0 (2019-09-14)

Initial chart release with Helm Operator [1.0.0-rc1](https://github.com/fluxcd/helm-operator/blob/master/CHANGELOG.md#100-rc1-2019-08-14)
