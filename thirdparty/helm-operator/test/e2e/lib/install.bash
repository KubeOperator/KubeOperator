#!/usr/bin/env bash

# shellcheck disable=SC1090
source "${E2E_DIR}/lib/defer.bash"
# shellcheck disable=SC1090
source "${E2E_DIR}/lib/template.bash"

function install_tiller() {
  if ! helm2 version > /dev/null 2>&1; then # only if helm isn't already installed
    kubectl --namespace "$E2E_NAMESPACE" create sa tiller
    kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount="$E2E_NAMESPACE":tiller
    helm2 init --tiller-namespace "$E2E_NAMESPACE" --service-account tiller --upgrade --wait
  fi
}

function uninstall_tiller() {
  # Note: helm reset --force will delete the Tiller
  # instance but will not delete release history.
  helm2 reset --tiller-namespace "$E2E_NAMESPACE" --force
  kubectl delete clusterrolebinding tiller-cluster-rule
  kubectl --namespace "$E2E_NAMESPACE" delete sa tiller
}

function install_helm_operator_with_helm() {
  local create_crds='true'
  if kubectl get crd helmreleases.helm.fluxcd.io > /dev/null 2>&1; then
    echo 'CRD existed, disabling CRD creation'
    create_crds='false'
  fi

  helm2 install --name helm-operator --wait \
    --tiller-namespace "${E2E_NAMESPACE}" \
    --namespace "${E2E_NAMESPACE}" \
    --set createCRD="${create_crds}" \
    --set chartsSyncInterval=3s \
    --set image.repository=docker.io/fluxcd/helm-operator \
    --set image.tag=latest \
    --set git.pollInterval=3s \
    --set git.config.secretName=gitconfig \
    --set git.config.enabled=true \
    --set-string git.config.data="${GITCONFIG}" \
    --set git.ssh.secretName=flux-git-deploy \
    --set-string git.ssh.known_hosts="${KNOWN_HOSTS}" \
    --set configureRepositories.enable=true \
    --set configureRepositories.repositories[0].name="stable" \
    --set configureRepositories.repositories[0].url="https://kubernetes-charts.storage.googleapis.com" \
    --set configureRepositories.repositories[1].name="podinfo" \
    --set configureRepositories.repositories[1].url="https://stefanprodan.github.io/podinfo" \
    --set helm.versions="${HELM_VERSION:-v2\,v3}" \
    --set tillerNamespace="${E2E_NAMESPACE}" \
    "${ROOT_DIR}/chart/helm-operator"
}

function uninstall_helm_operator_with_helm() {
  helm2 delete \
    --tiller-namespace "$E2E_NAMESPACE" \
    --purge helm-operator > /dev/null 2>&1

  kubectl delete crd helmreleases.helm.fluxcd.io > /dev/null 2>&1
}

function install_git_srv() {
  local external_access_result_var=${1}
  local kustomization_dir=${2:-base/gitsrv}
  local gen_dir
  gen_dir=$(mktemp -d)
  defer rm -rf "'$gen_dir'"
  ssh-keygen -t rsa -N "" -f "$gen_dir/id_rsa"
  kubectl create secret generic flux-git-deploy \
    --namespace="${E2E_NAMESPACE}" \
    --from-file="${FIXTURES_DIR}/known_hosts" \
    --from-file="$gen_dir/id_rsa" \
    --from-file=identity="$gen_dir/id_rsa" \
    --from-file="$gen_dir/id_rsa.pub"

  kubectl apply -n "${E2E_NAMESPACE}" -k "${FIXTURES_DIR}/kustom/${kustomization_dir}" >&3

  # Wait for the git server to be ready
  kubectl -n "${E2E_NAMESPACE}" rollout status deployment/gitsrv

  if [ -n "$external_access_result_var" ]; then
    local git_srv_podname
    git_srv_podname=$(kubectl get pod -n "${E2E_NAMESPACE}" -l name=gitsrv -o jsonpath="{['items'][0].metadata.name}")
    coproc kubectl port-forward -n "${E2E_NAMESPACE}" "$git_srv_podname" :22
    local local_port
    read -r local_port <&"${COPROC[0]}"-
    # shellcheck disable=SC2001
    local_port=$(echo "$local_port" | sed 's%.*:\([0-9]*\).*%\1%')
    local ssh_cmd="ssh -o UserKnownHostsFile=/dev/null  -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i $gen_dir/id_rsa -p $local_port"
    # return the ssh command needed for git, and the PID of the port-forwarding PID into a variable of choice
    eval "${external_access_result_var}=('$ssh_cmd' '$COPROC_PID')"
  fi
}

function uninstall_git_srv() {
  local kustomization_dir=${1:-base/gitsrv}

  # Silence secret deletion errors since the secret can be missing (deleted by uninstalling Flux)
  kubectl delete -n "${E2E_NAMESPACE}" secret flux-git-deploy &> /dev/null
  kubectl delete -n "${E2E_NAMESPACE}" -k "${FIXTURES_DIR}/kustom/${kustomization_dir}" >&3
}

function install_chartmuseum() {
  local external_access_result_var=${1}
  local kustomization_dir=${2:-base/chartmuseum}

  kubectl apply -n "${E2E_NAMESPACE}" -k "${FIXTURES_DIR}/kustom/${kustomization_dir}" >&3

  # Wait for the chartmuseum to become ready
  kubectl -n "${E2E_NAMESPACE}" rollout status deployment/chartmuseum

  if [ -n "$external_access_result_var" ]; then
    local chartmuseum_podname
    chartmuseum_podname=$(kubectl get pod -n "${E2E_NAMESPACE}" -l name=chartmuseum -o jsonpath="{['items'][0].metadata.name}")
    coproc kubectl port-forward -n "${E2E_NAMESPACE}" "$chartmuseum_podname" :8080
    local local_port
    read -r local_port <&"${COPROC[0]}"-
    # shellcheck disable=SC2001
    local_port=$(echo "$local_port" | sed 's%.*:\([0-9]*\).*%\1%')
    # return the ssh command needed for git, and the PID of the port-forwarding PID into a variable of choice
    eval "${external_access_result_var}=('$local_port' '$COPROC_PID')"
  fi
}

function uninstall_chartmuseum() {
  local kustomization_dir=${1:-base/chartmuseum}

  kubectl delete -n "${E2E_NAMESPACE}" -k "${FIXTURES_DIR}/kustom/${kustomization_dir}" >&3
}
