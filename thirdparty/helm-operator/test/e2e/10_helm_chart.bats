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

@test "Helm chart installation smoke test" {
  # The gitconfig secret must exist and have the right value
  poll_until_equals "gitconfig secret" "$GITCONFIG" "kubectl get secrets -n $E2E_NAMESPACE gitconfig -ojsonpath={..data.gitconfig} | base64 --decode"

  # Apply the HelmRelease fixtures
  kubectl apply -f "$FIXTURES_DIR/releases/git.yaml" >&3
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  poll_until_equals 'podinfo-helm-repository HelmRelease' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o 'custom-columns=status:status.releaseStatus' --no-headers"
  poll_until_equals 'podinfo-git HelmRelease' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o 'custom-columns=status:status.releaseStatus' --no-headers"
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
