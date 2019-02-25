#!/bin/bash
set -xeuo pipefail

# Essentially use a similar procedure other openshift-ansible PR tests use to
# determine which image tag should be used. This allows us to avoid hardcoding a
# specific version which quickly becomes stale.

if [ -n "${PAPR_BRANCH:-}" ]; then
  target_branch=$PAPR_BRANCH
else
  target_branch=$PAPR_PULL_TARGET_BRANCH
fi
target_branch_in=${target_branch}
if [[ "${target_branch}" =~ ^release- ]]; then
  target_branch="${target_branch/release-/}"
else
  dnf install -y sed
  target_branch="$( git describe | sed 's/^openshift-ansible-\([0-9]*\.[0-9]*\)\.[0-9]*-.*/\1/' )"
fi
export target_branch

# Need to define some git variables for rebase.
git config --global user.email "ci@openshift.org"
git config --global user.name "OpenShift Atomic CI"

# Rebase existing branch on the latest code locally, as PAPR running doesn't do merges
git fetch origin ${target_branch_in} && git rebase origin/${target_branch_in}

PAPR_INVENTORY=${PAPR_INVENTORY:-.papr.inventory}
PAPR_RUN_UPDATE=${PAPR_RUN_UPDATE:-0}
PAPR_UPGRADE_FROM=${PAPR_UPGRADE_FROM:-0}
PAPR_EXTRAVARS=""

# Replace current branch with PAPR_UPGRADE_FROM
if [[ "${PAPR_UPGRADE_FROM}" != "0" ]]; then
  git branch new-code
  git checkout release-${PAPR_UPGRADE_FROM}
  git clean -fdx
  PAPR_EXTRAVARS="-e openshift_release=${PAPR_UPGRADE_FROM}"
fi

pip install -r requirements.txt

# Human-readable output
export ANSIBLE_STDOUT_CALLBACK=debug

# ping the nodes to check they're responding and register their ostree versions
ansible -vv -i $PAPR_INVENTORY nodes -a 'rpm-ostree status'

# Make sure hostname -f returns correct node name
ansible -vv -i $PAPR_INVENTORY nodes -m setup
ansible -vv -i $PAPR_INVENTORY nodes -a "hostnamectl set-hostname {{ ansible_default_ipv4.address }}"
ansible -vv -i $PAPR_INVENTORY nodes -m setup -a "gather_subset=min"

upload_journals() {
  mkdir journals
  ansible -vvv -i $PAPR_INVENTORY all \
    -m shell -a 'journalctl --no-pager > /tmp/journal'
  ansible -vvv -i $PAPR_INVENTORY all \
    -m fetch -a "src=/tmp/journal dest=journals/{{ inventory_hostname }}.log flat=yes"

  # Split large files into parts, extracting a basename and preserving extention
  find . -iname "*.log" -execdir sh -c 'split -b 4m --numeric-suffixes --additional-suffix=.log {} $(basename {} .log)_' \; -execdir rm -rf {} \;
}

trap upload_journals ERR

# run the prerequisites play
ansible-playbook -vvv -i $PAPR_INVENTORY $PAPR_EXTRAVARS playbooks/prerequisites.yml

# run the actual installer
ansible-playbook -vvv -i $PAPR_INVENTORY $PAPR_EXTRAVARS playbooks/deploy_cluster.yml

# Restore the branch if needed
if [[ "${PAPR_UPGRADE_FROM}" != "0" ]]; then
  git checkout new-code
  git clean -fdx
  pip install -r requirements.txt
fi

# Run upgrade playbook
if [[ "${PAPR_RUN_UPDATE}" != "0" ]]; then
  update_version="$(echo $target_branch | sed 's/\./_/')"
  # Create basic node-group configmaps for upgrade
  ansible-playbook -vvv -i $PAPR_INVENTORY $PAPR_EXTRAVARS playbooks/openshift-master/openshift_node_group.yml
  ansible-playbook -vvv -i $PAPR_INVENTORY playbooks/byo/openshift-cluster/upgrades/v${update_version}/upgrade.yml | tee update.log
fi

upload_journals

### DISABLING TESTS FOR NOW, SEE:
### https://github.com/openshift/openshift-ansible/pull/6132

### # run a small subset of origin conformance tests to sanity
### # check the cluster NB: we run it on the master since we may
### # be in a different OSP network
### ssh ocp-master docker run --rm --net=host --privileged \
###   -v /etc/origin/master/admin.kubeconfig:/config \
###   registry.fedoraproject.org/fedora:27 sh -c \
###     '"dnf install -y origin-tests && \
###       KUBECONFIG=/config /usr/libexec/origin/extended.test --ginkgo.v=1 \
###         --ginkgo.noColor --ginkgo.focus=\"Services.*NodePort|EmptyDir\""'
