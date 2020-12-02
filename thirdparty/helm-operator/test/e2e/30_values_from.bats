#!/usr/bin/env bats

function setup() {
  # Load libraries in setup() to access BATS_* variables
  load lib/env
  load lib/defer
  load lib/install
  load lib/poll

  kubectl create namespace "$E2E_NAMESPACE"
  install_git_srv
  install_tiller
  install_helm_operator_with_helm
  kubectl create namespace "$DEMO_NAMESPACE"
}

@test "When valuesFrom.configMapKeyRefs are defined, they are sourced" {
  # Apply the HelmRelease fixtures
  kubectl apply -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/configmap.yaml" >&3

  # Wait for install failure
  poll_until_equals 'install failure due to missing config map' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/configmap-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmInstallFailed\")].message}'"

  # Add missing config map
  cat <<EOF | kubectl create -n "$DEMO_NAMESPACE" -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: podinfo-values
data:
  values.yaml: |
    replicaCount: 2
EOF

  # Wait for release recovery
  poll_until_equals 'recovery after adding config map' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/configmap-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Patch the release to make use of the config map in the other namespace
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/configmap.yaml" --type='json' -p="[{\"op\": \"add\", \"path\": \"/spec/valuesFrom/0/configMapKeyRef/namespace\", \"value\":\"$E2E_NAMESPACE\"}]" >&3

 # Wait for upgrade failure
  poll_until_equals 'upgrade failure due to missing config map in namespace' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/configmap-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmUpgradeFailed\")].message}'"

    # Add config map in different namespace
  cat <<EOF | kubectl create -n "$E2E_NAMESPACE" -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: podinfo-values
data:
  values.yaml: |
    replicaCount: 3
EOF

  # Wait for release recovery
  poll_until_equals 'recovery after adding config map in namespace' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/configmap-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"
}

@test "When valuesFrom.secretKeyRefs are defined, they are sourced" {
  # Apply the HelmRelease fixtures
  kubectl apply -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/secret.yaml" >&3

  # Wait for install failure
  poll_until_equals 'install failure due to missing secret' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/secret-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmInstallFailed\")].message}'"

  # Add missing config map
  cat <<EOF | kubectl create -n "$DEMO_NAMESPACE" -f -
apiVersion: v1
kind: Secret
metadata:
  name: podinfo-values
data:
  values.yaml: |
    cmVwbGljYUNvdW50OiAy
EOF

  # Wait for release recovery
  poll_until_equals 'recovery after adding secret' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/secret-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Patch the release to make use of the config map in the other namespace
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/secret.yaml" --type='json' -p="[{\"op\": \"add\", \"path\": \"/spec/valuesFrom/0/secretKeyRef/namespace\", \"value\":\"$E2E_NAMESPACE\"}]" >&3

 # Wait for upgrade failure
  poll_until_equals 'upgrade failure due to missing secret in namespace' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/secret-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmUpgradeFailed\")].message}'"

    # Add config map in different namespace
  cat <<EOF | kubectl create -n "$E2E_NAMESPACE" -f -
apiVersion: v1
kind: Secret
metadata:
  name: podinfo-values
data:
  values.yaml: |
    cmVwbGljYUNvdW50OiAz
EOF

  # Wait for release recovery
  poll_until_equals 'recovery after adding secret in namespace' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/secret-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"
}

@test "When valuesFrom.externalSourceRefs are defined, they are sourced" {
  # Apply the HelmRelease fixtures
  kubectl apply -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/externalsource.yaml" >&3

  # Wait for release
  poll_until_equals 'successful release with external source' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/external-source-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Patch the release with an invalid URL
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/externalsource.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/valuesFrom/0/externalSourceRef/url", "value": "https://raw.githubusercontent.com/stefanprodan/podinfo/3.2.0/charts/podinfo/invalid.yaml"}]' >&3

  # Wait for upgrade failure
  poll_until_equals 'upgrade failure due to invalid external source url' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/external-source-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmUpgradeFailed\")].message}'"

  # Make source optional
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/externalsource.yaml" --type='json' -p='[{"op": "add", "path": "/spec/valuesFrom/0/externalSourceRef/optional", "value": true}]' >&3

  # Wait for release recovery
  poll_until_equals 'recovery after making external source optional' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/external-source-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"
}

@test "When valuesFrom.chartFileRefs are defined, they are sourced" {
  # Apply the HelmRelease fixtures
  kubectl apply -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/chartfile.yaml" >&3

  # Wait for release
  poll_until_equals 'successful release with chart file' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/chartfile-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Patch release with invalid path
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/chartfile.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/valuesFrom/0/chartFileRef/path", "value": "invalid.yaml"}]' >&3

   # Wait for upgrade failure
  poll_until_equals 'upgrade failure due to invalid chart file' 'failed to compose values for chart release' "kubectl -n $DEMO_NAMESPACE get helmrelease/chartfile-values -o jsonpath='{.status.conditions[?(@.reason==\"HelmUpgradeFailed\")].message}'"

  # Make source optional
  kubectl patch -n "$DEMO_NAMESPACE" -f "$FIXTURES_DIR/values_from/chartfile.yaml" --type='json' -p='[{"op": "add", "path": "/spec/valuesFrom/0/chartFileRef/optional", "value": true}]' >&3

  # Wait for release recovery
  poll_until_equals 'recovery after making chart file optional' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/chartfile-values -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

}

function teardown() {
  run_deferred

  # Teardown is verbose when a test fails, and this will help most of the time
  # to determine _why_ it failed.
  kubectl logs -n "$E2E_NAMESPACE" deploy/helm-operator

  # Removing the operator also takes care of the global resources it installs.
  uninstall_helm_operator_with_helm
  uninstall_tiller
  # Removing the namespace also takes care of removing gitsrv.
  kubectl delete namespace "$E2E_NAMESPACE"
  # Only remove the demo workloads after the operator, so that they cannot be recreated.
  kubectl delete namespace "$DEMO_NAMESPACE"
}
