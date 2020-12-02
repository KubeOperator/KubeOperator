# Helm operator (`helm-operator`)

The Helm operator deals with Helm chart releases. The operator watches for
changes of Custom Resources of kind `HelmRelease`. It receives Kubernetes
Events and acts accordingly.

## Responsibilities

When the Helm Operator sees a `HelmRelease` resource in the
cluster, it either installs or upgrades the named Helm release so that
the chart is released as specified.

It will also notice when a `HelmRelease` resource is updated, and
take action accordingly.

## Setup and configuration

`helm-operator` requires setup and offers customization though a multitude of flags.

| Flag                        | Default                       | Purpose
| --------------------------  | ----------------------------- | ---
| `--log-format`              | `fmt`                         | Changes the logging format; `fmt` or `json`.
| `--workers`                 | `2`                           | Amount of workers processing releases.
| `--listen`                  | `:3030`                       | Listen address where `/metrics` and API will be served.
| **Cluster configuration**
| `--kubeconfig`              |                               | Path to a kubeconfig. Only required if out-of-cluster.
| `--master`                  |                               | The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
| `--allow-namespace`         |                               | If set, this limits the scope to a single namespace. if not specified, all namespaces will be watched.
| **Helm options**
| `--enabled-helm-versions`   | `v2,v3`                       | The Helm client versions supported by this operator instance
| `--helm-repository-import`  |                               | Targeted version and the path of the Helm repository index to import, i.e. `v3:/tmp/v3/index.yaml,v2:/tmp/v2/index.yaml`
| **Tiller options**
| `--tiller-ip`               |                               | Tiller IP address. Only required if out-of-cluster.
| `--tiller-port`             |                               | Tiller port.
| `--tiller-namespace`        |                               | Tiller namespace. If not provided, the default is kube-system.
| `--tiller-tls-enable`       | `false`                       | Enable TLS communication with Tiller. If provided, requires TLSKey and TLSCert to be provided as well.
| `--tiller-tls-verify`       | `false`                       | Verify TLS certificate from Tiller. Will enable TLS communication when provided.
| `--tiller-tls-key-path`     | `/etc/fluxd/helm/tls.key`     | Path to private key file used to communicate with the Tiller server.
| `--tiller-tls-cert-path`    | `/etc/fluxd/helm/tls.crt`     | Path to certificate file used to communicate with the Tiller server.
| `--tiller-tls-ca-cert-path` |                               | Path to CA certificate file used to validate the Tiller server. Required if tiller-tls-verify is enabled.
| `--tiller-tls-hostname`     |                               | The server name used to verify the hostname on the returned certificates from the Tiller server.
| **Reconciliation configuration**
| `--charts-sync-interval`    | `3m`                          | Period on which to reconcile the Helm releases with `HelmRelease` resources
| `--status-update-interval`  | `10s`                         | Period on which to update the Helm release status in `HelmRelease` resources
| `--log-release-diffs`       | `false`                       | Log the diff when a chart release diverges. **Potentially insecure due to logging of secret values.**
| **(Git sourced) chart changes** (none of these need overriding, usually)
| `--git-timeout`             | `20s`                         | Duration after which git operations time out.
| `--git-poll-interval`       | `5m`                          | Period on which to poll git chart sources for changes.
| `--update-chart-deps`       | `true`                        | Update chart dependencies before installing or upgrading a release.
