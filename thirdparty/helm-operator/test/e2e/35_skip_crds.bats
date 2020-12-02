#!/usr/bin/env bats

function setup() {
  # Load libraries in setup() to access BATS_* variables 
  load lib/env
  load lib/install
  load lib/poll

  kubectl create namespace "$E2E_NAMESPACE"
  install_git_srv
  install_tiller
  install_helm_operator_with_helm
  kubectl create namespace "$DEMO_NAMESPACE"
}

@test "When skipCRDs is set" {
  if [ "$HELM_VERSION" != "v3" ]; then
    skip
  fi

  # Apply the HelmRelease
  kubectl apply -f "$FIXTURES_DIR/releases/skip-crd.yaml" >&3
  poll_until_equals 'skip-crd HelmRelease' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/skip-crd -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Assert no CRDs were installed
  count=$(kubectl get crd --no-headers | grep 'konghq.com' | wc -l)
  [ "$count" -eq 0 ]
}

function teardown() {
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
