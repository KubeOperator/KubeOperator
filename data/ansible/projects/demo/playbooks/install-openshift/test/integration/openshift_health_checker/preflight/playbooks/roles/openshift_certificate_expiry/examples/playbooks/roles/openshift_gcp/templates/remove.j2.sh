#!/bin/bash

set -euo pipefail

function teardown_cmd() {
    a=( $@ )
    local name=$1
    a=( "${a[@]:1}" )
    local flag=0
    local found=
    for i in ${a[@]}; do
        if [[ "$i" == "--"* ]]; then
            found=true
            break
        fi
        flag=$((flag+1))
    done
    if [[ -z "${found}" ]]; then
      flag=$((flag+1))
    fi
    if gcloud --project "{{ openshift_gcp_project }}" ${a[@]::$flag} describe "${name}" ${a[@]:$flag} &>/dev/null; then
        gcloud --project "{{ openshift_gcp_project }}" ${a[@]::$flag} delete -q "${name}" ${a[@]:$flag}
    fi
}

function teardown() {
    for i in `seq 1 20`; do
        if teardown_cmd $@; then
            break
        fi
        sleep 0.5
    done
}

# Preemptively spin down the instances
{% for node_group in openshift_gcp_node_group_config %}
# scale down {{ node_group.name }}
(
    # performs a delete and scale down as one operation to ensure maximum parallelism
    if ! instances=$( gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed list-instances "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" --zone "{{ openshift_gcp_zone }}" --format='value[terminator=","](instance)' 2>/dev/null ); then
        exit 0
    fi
    instances="${instances%?}"
    if [[ -z "${instances}" ]]; then
        echo "warning: No instances in {{ node_group.name }}" 1>&2
        exit 0
    fi
    if ! gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed delete-instances "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" --zone "{{ openshift_gcp_zone }}" --instances "${instances}"; then
        echo "warning: Unable to scale down the node group {{ node_group.name }}" 1>&2
        exit 0
    fi
) &
{% endfor %}

# Bucket for registry
(
if gsutil ls -p "{{ openshift_gcp_project }}" "gs://{{ openshift_gcp_registry_bucket_name }}" &>/dev/null; then
    gsutil -m rm -r "gs://{{ openshift_gcp_registry_bucket_name }}"
fi
) &

# Project metadata prefixed with {{ openshift_gcp_prefix }}
(
    for key in $( gcloud --project "{{ openshift_gcp_project }}" compute project-info describe --flatten=commonInstanceMetadata.items[] '--format=value(commonInstanceMetadata.items.key)' ); do
        if [[ "${key}" == "{{ openshift_gcp_prefix }}"* ]]; then
            gcloud --project "{{ openshift_gcp_project }}" compute project-info remove-metadata "--keys=${key}"
        fi
    done
) &

# Instances and disks used for image building
(
    teardown "{{ openshift_gcp_prefix }}build-image-instance" compute instances --zone "{{ openshift_gcp_zone }}"
    teardown "{{ openshift_gcp_prefix }}build-image-instance" compute disks --zone "{{ openshift_gcp_zone }}"
) &

# DNS
(
dns_zone="{{ dns_managed_zone | default(openshift_gcp_prefix + 'managed-zone') }}"
if gcloud --project "{{ openshift_gcp_project }}" dns managed-zones describe "${dns_zone}" &>/dev/null; then
    # Retry DNS changes until they succeed since this may be a shared resource
    while true; do
        dns="${TMPDIR:-/tmp}/dns.yaml"
        rm -f "${dns}"

        # export all dns records that match into a zone format, and turn each line into a set of args for
        # record-sets transaction.
        gcloud dns record-sets export --project "{{ openshift_gcp_project }}" -z "${dns_zone}" --zone-file-format "${dns}"
        if grep -F -e '{{ openshift_master_cluster_hostname }}' -e '{{ openshift_master_cluster_public_hostname }}' -e '{{ wildcard_zone }}' "${dns}" | \
                awk '{ print "--name", $1, "--ttl", $2, "--type", $4, $5; }' > "${dns}.input"
        then
            rm -f "${dns}"
            gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns start -z "${dns_zone}"
            cat "${dns}.input" | xargs -L1 gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file="${dns}" remove -z "${dns_zone}"

            # Commit all DNS changes, retrying if preconditions are not met
            if ! out="$( gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns execute -z "${dns_zone}" 2>&1 )"; then
                rc=$?
                if [[ "${out}" == *"HTTPError 412: Precondition not met"* ]]; then
                    continue
                fi
                exit $rc
            fi
        fi
        rm "${dns}.input"
        break
    done
fi
) &

(
# Router network rules
teardown "{{ openshift_gcp_prefix }}router-network-lb-rule" compute forwarding-rules --region "{{ openshift_gcp_region }}"
teardown "{{ openshift_gcp_prefix }}router-network-lb-pool" compute target-pools --region "{{ openshift_gcp_region }}"
teardown "{{ openshift_gcp_prefix }}router-network-lb-health-check" compute http-health-checks
teardown "{{ openshift_gcp_prefix }}router-network-lb-ip" compute addresses --region "{{ openshift_gcp_region }}"

# Internal master network rules
teardown "{{ openshift_gcp_prefix }}master-network-lb-rule" compute forwarding-rules --region "{{ openshift_gcp_region }}"
teardown "{{ openshift_gcp_prefix }}master-network-lb-pool" compute target-pools --region "{{ openshift_gcp_region }}"
teardown "{{ openshift_gcp_prefix }}master-network-lb-health-check" compute http-health-checks
teardown "{{ openshift_gcp_prefix }}master-network-lb-ip" compute addresses --region "{{ openshift_gcp_region }}"
) &

(
# Master SSL network rules
teardown "{{ openshift_gcp_prefix }}master-ssl-lb-rule" compute forwarding-rules --global
teardown "{{ openshift_gcp_prefix }}master-ssl-lb-target" compute target-tcp-proxies
teardown "{{ openshift_gcp_prefix }}master-ssl-lb-ip" compute addresses --global
teardown "{{ openshift_gcp_prefix }}master-ssl-lb-backend" compute backend-services --global
teardown "{{ openshift_gcp_prefix }}master-ssl-lb-health-check" compute health-checks
) &

#Firewall rules
#['name']='parameters for "gcloud compute firewall-rules create"'
#For all possible parameters see: gcloud compute firewall-rules create --help
declare -A FW_RULES=(
  ['icmp']=""
  ['ssh-external']=""
  ['ssh-internal']=""
  ['master-internal']=""
  ['master-external']=""
  ['node-internal']=""
  ['infra-node-internal']=""
  ['infra-node-external']=""
)
for rule in "${!FW_RULES[@]}"; do
    ( if gcloud --project "{{ openshift_gcp_project }}" compute firewall-rules describe "{{ openshift_gcp_prefix }}$rule" &>/dev/null; then
        # retry a few times because this call can be flaky
        for i in `seq 1 3`; do 
            if gcloud -q --project "{{ openshift_gcp_project }}" compute firewall-rules delete "{{ openshift_gcp_prefix }}$rule"; then
                break
            fi
        done
    fi ) &
done

for i in `jobs -p`; do wait $i; done

{% for node_group in openshift_gcp_node_group_config %}
# teardown {{ node_group.name }} - any load balancers referencing these groups must be removed
(
    teardown "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" compute instance-groups managed --zone "{{ openshift_gcp_zone }}"
    teardown "{{ openshift_gcp_prefix }}instance-template-{{ node_group.name }}" compute instance-templates
) &
{% endfor %}

for i in `jobs -p`; do wait $i; done

# Images specifically located under this cluster prefix family
for name in $( gcloud --project "{{ openshift_gcp_project }}" compute images list "--filter=family={{ openshift_gcp_prefix }}images" '--format=value(name)' ); do
    ( gcloud --project "{{ openshift_gcp_project }}" compute images delete "${name}" ) &
done

# Network
( teardown "{{ openshift_gcp_network_name }}" compute networks ) &

for i in `jobs -p`; do wait $i; done