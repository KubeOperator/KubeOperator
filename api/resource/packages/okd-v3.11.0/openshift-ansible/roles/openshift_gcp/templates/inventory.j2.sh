#!/bin/sh

export GCE_PROJECT="{{ openshift_gcp_project }}"
export GCE_ZONE="{{ openshift_gcp_zone }}"
export GCE_EMAIL="{{ (lookup('file', openshift_gcp_iam_service_account_keyfile ) | from_json ).client_email }}"
export GCE_PEM_FILE_PATH="/tmp/gce.pem"
export INVENTORY_IP_TYPE="{{ inventory_ip_type }}"
export GCE_TAGGED_INSTANCES="{{ openshift_gcp_prefix }}ocp"