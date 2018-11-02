#!/bin/bash
#
# This script is a startup script for bootstrapping a GCP node
# from a config stored in the project metadata. It loops until
# it finds the script and then starts the origin-node service.
# TODO: generalize

set -o errexit
set -o nounset
set -o pipefail

if [[ "$( curl "http://metadata.google.internal/computeMetadata/v1/instance/attributes/bootstrap" -H "Metadata-Flavor: Google" )" != "true" ]]; then
  echo "info: Bootstrap is not enabled for this instance, skipping" 1>&2
  exit 0
fi

if ! id=$( curl "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-id" -H "Metadata-Flavor: Google" ); then
  echo "error: Unable to get cluster-id for instance from cluster metadata" 1>&2
  exit 1
fi

if ! node_group=$( curl "http://metadata.google.internal/computeMetadata/v1/instance/attributes/node-group" -H "Metadata-Flavor: Google" ); then
  echo "error: Unable to get node-group for instance from cluster metadata" 1>&2
  exit 1
fi

if ! config=$( curl -f "http://metadata.google.internal/computeMetadata/v1/instance/attributes/bootstrap-config" -H "Metadata-Flavor: Google" 2>/dev/null ); then
  while true; do
    if config=$( curl -f "http://metadata.google.internal/computeMetadata/v1/project/attributes/${id}-bootstrap-config" -H "Metadata-Flavor: Google" 2>/dev/null ); then
      break
    fi
    echo "info: waiting for ${id}-bootstrap-config to become available in cluster metadata ..." 1>&2
    sleep 5
  done
fi

echo "Got bootstrap config from metadata"
mkdir -p /etc/origin/node
echo -n "${config}" > /etc/origin/node/bootstrap.kubeconfig
echo "BOOTSTRAP_CONFIG_NAME=node-config-${node_group}" >> /etc/sysconfig/origin-node
systemctl enable origin-node
systemctl start origin-node
