#!/bin/bash

set -euo pipefail

metadata=""
if [[ -n "{{ openshift_gcp_startup_script_file }}" ]]; then
    if [[ ! -f "{{ openshift_gcp_startup_script_file }}" ]]; then
        echo "Startup script file missing at {{ openshift_gcp_startup_script_file }} from=$(pwd)"
        exit 1
    fi
    metadata+="--metadata-from-file=startup-script={{ openshift_gcp_startup_script_file }}"
fi
if [[ -n "{{ openshift_gcp_user_data_file }}" ]]; then
    if [[ ! -f "{{ openshift_gcp_user_data_file }}" ]]; then
        echo "User data file missing at {{ openshift_gcp_user_data_file }}"
        exit 1
    fi
    if [[ -n "${metadata}" ]]; then
        metadata+=","
    else
        metadata="--metadata-from-file="
    fi
    metadata+="user-data={{ openshift_gcp_user_data_file }}"
fi

# Select image or image family
image="{{ openshift_gcp_image }}"
if ! gcloud --project "{{ openshift_gcp_project }}" compute images describe "${image}" &>/dev/null; then
    if ! gcloud --project "{{ openshift_gcp_project }}" compute images describe-from-family "${image}" &>/dev/null; then
        echo "No compute image or image-family found, create an image named '{{ openshift_gcp_image }}' to continue'"
        exit 1
    fi
    image="family/${image}"
fi

### PROVISION THE INFRASTRUCTURE ###

dns_zone="{{ dns_managed_zone | default(openshift_gcp_prefix + 'managed-zone') }}"

# Check the DNS managed zone in Google Cloud DNS, create it if it doesn't exist and exit after printing NS servers
if ! gcloud --project "{{ openshift_gcp_project }}" dns managed-zones describe "${dns_zone}" &>/dev/null; then
    echo "DNS zone '${dns_zone}' doesn't exist. Must be configured prior to running this script"
    exit 1
fi

# Create network
if ! gcloud --project "{{ openshift_gcp_project }}" compute networks describe "{{ openshift_gcp_network_name }}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute networks create "{{ openshift_gcp_network_name }}" --mode "auto"
else
    echo "Network '{{ openshift_gcp_network_name }}' already exists"
fi

# Firewall rules in a form:
# ['name']='parameters for "gcloud compute firewall-rules create"'
# For all possible parameters see: gcloud compute firewall-rules create --help
range=""
if [[ -n "{{ openshift_node_port_range }}" ]]; then
    range=",tcp:{{ openshift_node_port_range }},udp:{{ openshift_node_port_range }}"
fi
declare -A FW_RULES=(
  ['icmp']='--allow icmp'
  ['ssh-external']='--allow tcp:22'
  ['ssh-internal']='--allow tcp:22 --source-tags bastion'
  ['master-internal']="--allow tcp:2224,tcp:2379,tcp:2380,tcp:4001,udp:4789,udp:5404,udp:5405,tcp:8053,udp:8053,tcp:8444,tcp:10250,tcp:10255,udp:10255,tcp:24224,udp:24224 --source-tags ocp --target-tags ocp-master"
  ['master-external']="--allow tcp:80,tcp:443,tcp:1936,tcp:8080,tcp:8443${range} --target-tags ocp-master"
  ['node-internal']="--allow udp:4789,tcp:10250,tcp:10255,udp:10255,tcp:9000-10000 --source-tags ocp --target-tags ocp-node,ocp-infra-node"
  ['infra-node-internal']="--allow tcp:5000 --source-tags ocp --target-tags ocp-infra-node"
  ['infra-node-external']="--allow tcp:80,tcp:443,tcp:1936${range} --target-tags ocp-infra-node"
)
for rule in "${!FW_RULES[@]}"; do
    ( if ! gcloud --project "{{ openshift_gcp_project }}" compute firewall-rules describe "{{ openshift_gcp_prefix }}$rule" &>/dev/null; then
        gcloud --project "{{ openshift_gcp_project }}" compute firewall-rules create "{{ openshift_gcp_prefix }}$rule" --network "{{ openshift_gcp_network_name }}" ${FW_RULES[$rule]}
    else
        echo "Firewall rule '{{ openshift_gcp_prefix }}${rule}' already exists"
    fi ) &
done


# Master IP
( if ! gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-ssl-lb-ip" --global &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute addresses create "{{ openshift_gcp_prefix }}master-ssl-lb-ip" --global
else
    echo "IP '{{ openshift_gcp_prefix }}master-ssl-lb-ip' already exists"
fi ) &

# Internal master IP
( if ! gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-network-lb-ip" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute addresses create "{{ openshift_gcp_prefix }}master-network-lb-ip" --region "{{ openshift_gcp_region }}"
else
    echo "IP '{{ openshift_gcp_prefix }}master-network-lb-ip' already exists"
fi ) &

# Router IP
( if ! gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}router-network-lb-ip" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute addresses create "{{ openshift_gcp_prefix }}router-network-lb-ip" --region "{{ openshift_gcp_region }}"
else
    echo "IP '{{ openshift_gcp_prefix }}router-network-lb-ip' already exists"
fi ) &


{% for node_group in openshift_gcp_node_group_config %}
# configure {{ node_group.name }}
(
    if ! gcloud --project "{{ openshift_gcp_project }}" compute instance-templates describe "{{ openshift_gcp_prefix }}instance-template-{{ node_group.name }}" &>/dev/null; then
        gcloud --project "{{ openshift_gcp_project }}" compute instance-templates create "{{ openshift_gcp_prefix }}instance-template-{{ node_group.name }}" \
                --machine-type "{{ node_group.machine_type }}" --network "{{ openshift_gcp_network_name }}" \
                --tags "{{ openshift_gcp_prefix }}ocp,ocp,ocp-bootstrap,{{ node_group.tags }}" \
                --boot-disk-size "{{ node_group.boot_disk_size }}" --boot-disk-type "pd-ssd" \
                --scopes "logging-write,monitoring-write,useraccounts-ro,service-control,service-management,storage-ro,compute-rw" \
                --image "{{ node_group.image | default('${image}') }}" ${metadata}  \
                --metadata "bootstrap={{ node_group.bootstrap | default(False) | bool | to_json }},cluster-id={{ openshift_gcp_prefix + openshift_gcp_clusterid }},node-group={{ node_group.name }}"
    else
        echo "Instance template '{{ openshift_gcp_prefix }}instance-template-{{ node_group.name }}' already exists"
    fi

    # Create instance group
    if ! gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed describe "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" --zone "{{ openshift_gcp_zone }}" &>/dev/null; then
        gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed create "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" \
                --zone "{{ openshift_gcp_zone }}" --template "{{ openshift_gcp_prefix }}instance-template-{{ node_group.name }}" --size "{{ node_group.scale }}"
    else
        echo "Instance group '{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}' already exists"
    fi
) &
{% endfor %}

for i in `jobs -p`; do wait $i; done


# Configure the master external LB rules
(
# Master health check
if ! gcloud --project "{{ openshift_gcp_project }}" compute health-checks describe "{{ openshift_gcp_prefix }}master-ssl-lb-health-check" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute health-checks create https "{{ openshift_gcp_prefix }}master-ssl-lb-health-check" --port "{{ internal_console_port }}" --request-path "/healthz"
else
    echo "Health check '{{ openshift_gcp_prefix }}master-ssl-lb-health-check' already exists"
fi

gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed set-named-ports "{{ openshift_gcp_prefix }}ig-m" \
        --zone "{{ openshift_gcp_zone }}" --named-ports "{{ openshift_gcp_prefix }}port-name-master:{{ internal_console_port }}"

# Master backend service
if ! gcloud --project "{{ openshift_gcp_project }}" compute backend-services describe "{{ openshift_gcp_prefix }}master-ssl-lb-backend" --global &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute backend-services create "{{ openshift_gcp_prefix }}master-ssl-lb-backend" --health-checks "{{ openshift_gcp_prefix }}master-ssl-lb-health-check" --port-name "{{ openshift_gcp_prefix }}port-name-master" --protocol "TCP" --global --timeout="{{ openshift_gcp_master_lb_timeout }}"
    gcloud --project "{{ openshift_gcp_project }}" compute backend-services add-backend "{{ openshift_gcp_prefix }}master-ssl-lb-backend" --instance-group "{{ openshift_gcp_prefix }}ig-m" --global --instance-group-zone "{{ openshift_gcp_zone }}"
else
    echo "Backend service '{{ openshift_gcp_prefix }}master-ssl-lb-backend' already exists"
fi

# Master tcp proxy target
if ! gcloud --project "{{ openshift_gcp_project }}" compute target-tcp-proxies describe "{{ openshift_gcp_prefix }}master-ssl-lb-target" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute target-tcp-proxies create "{{ openshift_gcp_prefix }}master-ssl-lb-target" --backend-service "{{ openshift_gcp_prefix }}master-ssl-lb-backend"
else
    echo "Proxy target '{{ openshift_gcp_prefix }}master-ssl-lb-target' already exists"
fi

# Master forwarding rule
if ! gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules describe "{{ openshift_gcp_prefix }}master-ssl-lb-rule" --global &>/dev/null; then
    IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-ssl-lb-ip" --global --format='value(address)')
    gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules create "{{ openshift_gcp_prefix }}master-ssl-lb-rule" --address "$IP" --global --ports "{{ console_port }}" --target-tcp-proxy "{{ openshift_gcp_prefix }}master-ssl-lb-target"
else
    echo "Forwarding rule '{{ openshift_gcp_prefix }}master-ssl-lb-rule' already exists"
fi
) &


# Configure the master internal LB rules
(
# Internal master health check
if ! gcloud --project "{{ openshift_gcp_project }}" compute http-health-checks describe "{{ openshift_gcp_prefix }}master-network-lb-health-check" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute http-health-checks create "{{ openshift_gcp_prefix }}master-network-lb-health-check" --port "8080" --request-path "/healthz"
else
    echo "Health check '{{ openshift_gcp_prefix }}master-network-lb-health-check' already exists"
fi

# Internal master target pool
if ! gcloud --project "{{ openshift_gcp_project }}" compute target-pools describe "{{ openshift_gcp_prefix }}master-network-lb-pool" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute target-pools create "{{ openshift_gcp_prefix }}master-network-lb-pool" --http-health-check "{{ openshift_gcp_prefix }}master-network-lb-health-check" --region "{{ openshift_gcp_region }}"
else
    echo "Target pool '{{ openshift_gcp_prefix }}master-network-lb-pool' already exists"
fi

# Internal master forwarding rule
if ! gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules describe "{{ openshift_gcp_prefix }}master-network-lb-rule" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-network-lb-ip" --region "{{ openshift_gcp_region }}" --format='value(address)')
    gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules create "{{ openshift_gcp_prefix }}master-network-lb-rule" --address "$IP" --region "{{ openshift_gcp_region }}" --target-pool "{{ openshift_gcp_prefix }}master-network-lb-pool"
else
    echo "Forwarding rule '{{ openshift_gcp_prefix }}master-network-lb-rule' already exists"
fi
) &


# Configure the infra node rules
(
# Router health check
if ! gcloud --project "{{ openshift_gcp_project }}" compute http-health-checks describe "{{ openshift_gcp_prefix }}router-network-lb-health-check" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute http-health-checks create "{{ openshift_gcp_prefix }}router-network-lb-health-check" --port "1936" --request-path "/healthz"
else
    echo "Health check '{{ openshift_gcp_prefix }}router-network-lb-health-check' already exists"
fi

# Router target pool
if ! gcloud --project "{{ openshift_gcp_project }}" compute target-pools describe "{{ openshift_gcp_prefix }}router-network-lb-pool" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    gcloud --project "{{ openshift_gcp_project }}" compute target-pools create "{{ openshift_gcp_prefix }}router-network-lb-pool" --http-health-check "{{ openshift_gcp_prefix }}router-network-lb-health-check" --region "{{ openshift_gcp_region }}"
else
    echo "Target pool '{{ openshift_gcp_prefix }}router-network-lb-pool' already exists"
fi

# Router forwarding rule
if ! gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules describe "{{ openshift_gcp_prefix }}router-network-lb-rule" --region "{{ openshift_gcp_region }}" &>/dev/null; then
    IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}router-network-lb-ip" --region "{{ openshift_gcp_region }}" --format='value(address)')
    gcloud --project "{{ openshift_gcp_project }}" compute forwarding-rules create "{{ openshift_gcp_prefix }}router-network-lb-rule" --address "$IP" --region "{{ openshift_gcp_region }}" --target-pool "{{ openshift_gcp_prefix }}router-network-lb-pool"
else
    echo "Forwarding rule '{{ openshift_gcp_prefix }}router-network-lb-rule' already exists"
fi
) &

for i in `jobs -p`; do wait $i; done

# set the target pools
(
if [[ "ig-m" == "{{ openshift_gcp_infra_network_instance_group }}" ]]; then
    gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed set-target-pools "{{ openshift_gcp_prefix }}ig-m" --target-pools "{{ openshift_gcp_prefix }}master-network-lb-pool,{{ openshift_gcp_prefix }}router-network-lb-pool" --zone "{{ openshift_gcp_zone }}"
else
    gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed set-target-pools "{{ openshift_gcp_prefix }}ig-m" --target-pools "{{ openshift_gcp_prefix }}master-network-lb-pool" --zone "{{ openshift_gcp_zone }}"
    gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed set-target-pools "{{ openshift_gcp_prefix }}{{ openshift_gcp_infra_network_instance_group }}" --target-pools "{{ openshift_gcp_prefix }}router-network-lb-pool" --zone "{{ openshift_gcp_zone }}"
fi
) &

# configure DNS
(
# Retry DNS changes until they succeed since this may be a shared resource
while true; do
    dns="${TMPDIR:-/tmp}/dns.yaml"
    rm -f $dns

    # DNS record for master lb
    if ! gcloud --project "{{ openshift_gcp_project }}" dns record-sets list -z "${dns_zone}" --name "{{ openshift_master_cluster_public_hostname }}" 2>/dev/null | grep -q "{{ openshift_master_cluster_public_hostname }}"; then
        IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-ssl-lb-ip" --global --format='value(address)')
        if [[ ! -f $dns ]]; then
            gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns start -z "${dns_zone}"
        fi
        gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns add -z "${dns_zone}" --ttl {{ openshift_gcp_master_dns_ttl }} --name "{{ openshift_master_cluster_public_hostname }}." --type A "$IP"
    else
        echo "DNS record for '{{ openshift_master_cluster_public_hostname }}' already exists"
    fi

    # DNS record for internal master lb
    if ! gcloud --project "{{ openshift_gcp_project }}" dns record-sets list -z "${dns_zone}" --name "{{ openshift_master_cluster_hostname }}" 2>/dev/null | grep -q "{{ openshift_master_cluster_hostname }}"; then
        IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}master-network-lb-ip" --region "{{ openshift_gcp_region }}" --format='value(address)')
        if [[ ! -f $dns ]]; then
            gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns start -z "${dns_zone}"
        fi
        gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns add -z "${dns_zone}" --ttl {{ openshift_gcp_master_dns_ttl }} --name "{{ openshift_master_cluster_hostname }}." --type A "$IP"
    else
        echo "DNS record for '{{ openshift_master_cluster_hostname }}' already exists"
    fi

    # DNS record for router lb
    if ! gcloud --project "{{ openshift_gcp_project }}" dns record-sets list -z "${dns_zone}" --name "{{ wildcard_zone }}" 2>/dev/null | grep -q "{{ wildcard_zone }}"; then
        IP=$(gcloud --project "{{ openshift_gcp_project }}" compute addresses describe "{{ openshift_gcp_prefix }}router-network-lb-ip" --region "{{ openshift_gcp_region }}" --format='value(address)')
        if [[ ! -f $dns ]]; then
            gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns start -z "${dns_zone}"
        fi
        gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns add -z "${dns_zone}" --ttl {{ openshift_gcp_master_dns_ttl }} --name "{{ wildcard_zone }}." --type A "$IP"
        gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns add -z "${dns_zone}" --ttl {{ openshift_gcp_master_dns_ttl }} --name "*.{{ wildcard_zone }}." --type CNAME "{{ wildcard_zone }}."
    else
        echo "DNS record for '{{ wildcard_zone }}' already exists"
    fi

    # Commit all DNS changes, retrying if preconditions are not met
    if [[ -f $dns ]]; then
        if ! out="$( gcloud --project "{{ openshift_gcp_project }}" dns record-sets transaction --transaction-file=$dns execute -z "${dns_zone}" 2>&1 )"; then
            rc=$?
            if [[ "${out}" == *"HTTPError 412: Precondition not met"* ]]; then
                continue
            fi
            exit $rc
        fi
    fi
    break
done
) &

# Create bucket for registry
(
if ! gsutil ls -p "{{ openshift_gcp_project }}" "gs://{{ openshift_gcp_registry_bucket_name }}" &>/dev/null; then
    gsutil mb -p "{{ openshift_gcp_project }}" -l "{{ openshift_gcp_region }}" "gs://{{ openshift_gcp_registry_bucket_name }}"
else
    echo "Bucket '{{ openshift_gcp_registry_bucket_name }}' already exists"
fi
) &

# wait until all node groups are stable
{% for node_group in openshift_gcp_node_group_config %}
{% if node_group.wait_for_stable | default(False) or not (node_group.bootstrap | default(False)) %}
# wait for stable {{ node_group.name }}
( gcloud --project "{{ openshift_gcp_project }}" compute instance-groups managed wait-until-stable "{{ openshift_gcp_prefix }}ig-{{ node_group.suffix }}" --zone "{{ openshift_gcp_zone }}" --timeout=600 ) &
{% else %}
# not waiting for {{ node_group.name }} due to bootstrapping
{% endif %}
{% endfor %}


for i in `jobs -p`; do wait $i; done
