#!/usr/bin/env bats

function setup() {
  # Load libraries in setup() to access BATS_* variables
  load lib/env
  load lib/defer
  load lib/helm
  load lib/install
  load lib/poll

  kubectl create namespace "$E2E_NAMESPACE"

  # Install the git server, allowing external access
  install_git_srv git_srv_result
  # shellcheck disable=SC2154
  export GIT_SSH_COMMAND="${git_srv_result[0]}"
  # Teardown the created port-forward to gitsrv.
  defer kill "${git_srv_result[1]}"

  install_tiller
  install_helm_operator_with_helm

  kubectl create namespace "$DEMO_NAMESPACE"
}

@test "When a mutation it Git is made, a release is upgraded" {
  # Apply the HelmRelease fixture
  kubectl apply -f "$FIXTURES_DIR/releases/git.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release to be deployed' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Clone the charts repository
  local clone_dir
  clone_dir="$(mktemp -d)"
  defer rm -rf "'$clone_dir'"
  git clone -b master ssh://git@localhost/git-server/repos/cluster.git "$clone_dir"
  cd "$clone_dir"

  # Make a chart template mutation in Git without bumping the version number
  sed -i 's%these commands:%these commands;%' charts/podinfo/templates/NOTES.txt
  git add charts/podinfo/templates/NOTES.txt
  git -c 'user.email=foo@bar.com' -c 'user.name=Foo' commit -m "Modify NOTES.txt"

  # Record new HEAD and push change
  head_hash=$(git rev-list -n 1 HEAD)
  git push >&3

  # Assert change is rolled out
  poll_until_equals 'release chart update' "successfully cloned chart revision: $head_hash" "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o jsonpath='{.status.conditions[?(@.type==\"ChartFetched\")].message}'"
  poll_until_equals 'revision match' "$head_hash" "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o jsonpath='{.status.revision}'"
}

@test "When a values.yaml change in Git is made, a release is upgraded" {
  # Apply the HelmRelease fixture
  kubectl apply -f "$FIXTURES_DIR/releases/git.yaml" >&3

  # Wait for it to be deployed
  poll_until_equals 'release to be deployed' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o 'custom-columns=status:status.releaseStatus' --no-headers"

  # Clone the charts repository
  local clone_dir
  clone_dir="$(mktemp -d)"
  defer rm -rf "'$clone_dir'"
  git clone -b master ssh://git@localhost/git-server/repos/cluster.git "$clone_dir"
  cd "$clone_dir"

  # Make a values.yaml mutation in Git
  sed -i 's%replicaCount: 1%replicaCount: 2%' charts/podinfo/values.yaml
  git add charts/podinfo/values.yaml
  git -c 'user.email=foo@bar.com' -c 'user.name=Foo' commit -m "Change replicaCount to 2"

  # Record new HEAD and push change
  head_hash=$(git rev-list -n 1 HEAD)
  git push >&3

  # Assert change is rolled out
  poll_until_equals 'release chart update' "successfully cloned chart revision: $head_hash" "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o jsonpath='{.status.conditions[?(@.type==\"ChartFetched\")].message}'"
  poll_until_equals 'revision match' "$head_hash" "kubectl -n $DEMO_NAMESPACE get helmrelease/podinfo-git -o jsonpath='{.status.revision}'"
}

@test "When a HelmRelease is nested in a chart, an upgrade does succeed" {
  # Install chartmuseum
  install_chartmuseum chartmuseum_result
  # shellcheck disable=SC2154
  local CHARTMUSEUM_URL="http://localhost:${chartmuseum_result[0]}"
    # Teardown the created port-forward to chartmusem.
  defer kill "${chartmuseum_result[1]}"

  # Package and upload chart fixture
  package_and_upload_chart "$FIXTURES_DIR/charts/nested-helmrelease" "$CHARTMUSEUM_URL"

  # Apply the HelmRelease fixture
  kubectl apply -f "$FIXTURES_DIR/releases/nested-helmrelease.yaml" >&3

  # Wait for it and the child release to be deployed
  poll_until_equals 'release to be deployed' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/nested-helmrelease -o 'custom-columns=status:status.releaseStatus' --no-headers"
  poll_until_equals 'child release to be deployed' 'deployed' "kubectl -n $DEMO_NAMESPACE get helmrelease/nested-helmrelease-child -o 'custom-columns=status:status.releaseStatus' --no-headers"

  releaseGen=$(kubectl -n "$DEMO_NAMESPACE" get helmrelease/nested-helmrelease -o 'custom-columns=status:status.observedGeneration' --no-headers)
  childReleaseGen=$(kubectl -n "$DEMO_NAMESPACE" get helmrelease/nested-helmrelease -o 'custom-columns=status:status.observedGeneration' --no-headers)

  # Patch release
  kubectl patch -f "$FIXTURES_DIR/releases/nested-helmrelease.yaml" --type='json' -p='[{"op": "replace", "path": "/spec/values/nested/deeper/deepest/image/tag", "value": "1.1.0"}]' >&3

  # Wait for patch to be processed
  poll_until_equals 'patch to be processed for release' "$((releaseGen+1))" "kubectl -n $DEMO_NAMESPACE get helmrelease/nested-helmrelease -o 'custom-columns=status:status.observedGeneration' --no-headers"
  poll_until_equals 'patch to be processed for child release' "$((childReleaseGen+1))" "kubectl -n $DEMO_NAMESPACE get helmrelease/nested-helmrelease -o 'custom-columns=status:status.observedGeneration' --no-headers"

  # Assert successful release
  releaseStatus=$(kubectl -n "$DEMO_NAMESPACE" get helmrelease/nested-helmrelease -o jsonpath='{.status.conditions[?(@.type=="Released")].reason}')
  [ "$releaseStatus" = "HelmSuccess" ]
  childReleaseStatus=$(kubectl -n "$DEMO_NAMESPACE" get helmrelease/nested-helmrelease-child -o jsonpath='{.status.conditions[?(@.type=="Released")].reason}')
  [ "$childReleaseStatus" = "HelmSuccess" ]
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
