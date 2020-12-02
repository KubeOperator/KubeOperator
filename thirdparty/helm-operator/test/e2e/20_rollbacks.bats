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

@test "When rollback.enable is set, failed releases are rolled back" {
  # Apply the HelmRelease
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release deploy' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Apply a faulty patch
  kubectl patch -f "$FIXTURES_DIR/releases/helm-repository.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/values/faults/unready", "value":"true"}]' >&3

  # Wait for release failure
  poll_until_equals 'upgrade failure' 'HelmUpgradeFailed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Wait for rollback
  poll_until_equals 'rollback' 'True' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.conditions[?(@.type==\"RolledBack\")].status}'"

  # Apply fix patch
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  # Assert recovery
  poll_until_equals 'recovery' 'HelmSuccess' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"
}

@test "When rollback.retry is set, upgrades are reattempted after a rollback" {
  # Apply the HelmRelease
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release deploy' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Apply a faulty patch and enable retries
  kubectl patch -f "$FIXTURES_DIR/releases/helm-repository.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/values/faults/unready", "value": true},{"op": "add", "path": "/spec/rollback/retry", "value": true}]' >&3

  # Wait for release failure
  poll_until_equals 'upgrade failure' 'HelmUpgradeFailed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Wait for rollback count to increase
  poll_until_equals 'rollback count == 3' '3' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.rollbackCount}'"

  # Apply fix patch
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  # Assert rollback count is reset
  poll_until_equals 'rollback count reset' '' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.rollbackCount}'" >&3
}

@test "When rollback.maxRetries is set to 1,  upgrade is only retried once" {
  # Apply the HelmRelease
  kubectl apply -f "$FIXTURES_DIR/releases/helm-repository.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release deploy' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Apply a faulty patch and enable retries
  kubectl patch -f "$FIXTURES_DIR/releases/helm-repository.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/values/faults/unready", "value": true},{"op": "add", "path": "/spec/rollback/retry", "value": true},{"op": "add", "path": "/spec/rollback/maxRetries", "value": 1}]' >&3

  # Wait for release failure
  poll_until_equals 'upgrade failure' 'HelmUpgradeFailed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.conditions[?(@.type==\"Released\")].reason}'"

  # Wait for rollback count to increase
  poll_until_equals 'rollback count == 2' '2' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-helm-repository -o jsonpath='{.status.rollbackCount}'"

  # Wait for dry-run to be compared, instead of retry
  poll_until_true 'dry-run comparison to failed release' "kubectl -n $E2E_NAMESPACE logs deploy/helm-operator | grep -E \"comparing dry-run output with latest failed release\""
}

@test "When rollback.enable is set, validation error does not trigger a rollback" {
  if [ "$HELM_VERSION" != "v3" ]; then
    skip
  fi

  # Apply the HelmRelease
  kubectl apply -f "$FIXTURES_DIR/releases/git.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release deploy' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Apply a faulty patch
  kubectl patch -f "$FIXTURES_DIR/releases/git.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/values/replicaCount", "value":"faulty"}]' >&3

  # Wait for release failure
  poll_until_equals 'upgrade failure' 'False' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o jsonpath='{.status.conditions[?(@.reason==\"HelmUpgradeFailed\")].status}'"

  # Assert release version
  version=$(kubectl exec -n "$E2E_NAMESPACE" deploy/helm-operator -- helm3 status podinfo-git --namespace "$DEMO_NAMESPACE" -o json | jq .version)
  [ "$version" -eq 1 ]

  # Assert rollback count is zero
  count=$(kubectl -n "$DEMO_NAMESPACE" get helmrelease/podinfo-git -o jsonpath='{.status.rollbackCount}')
  [ -z "$count" ]
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
