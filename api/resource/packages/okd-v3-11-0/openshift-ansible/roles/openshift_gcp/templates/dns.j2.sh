#!/bin/bash

set -euo pipefail

dns_zone="{{ dns_managed_zone | default(openshift_gcp_prefix + 'managed-zone') }}"

# Check the DNS managed zone in Google Cloud DNS, create it if it doesn't exist
if ! gcloud --project "{{ openshift_gcp_project }}" dns managed-zones describe "${dns_zone}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" dns managed-zones create "${dns_zone}" --dns-name "{{ public_hosted_zone }}" --description "{{ public_hosted_zone }} domain" >/dev/null
fi

# Always output the expected nameservers as a comma delimited list
gcloud --project "{{ openshift_gcp_project }}" dns managed-zones describe "${dns_zone}" --format='value(nameServers)' | tr ';' ','
