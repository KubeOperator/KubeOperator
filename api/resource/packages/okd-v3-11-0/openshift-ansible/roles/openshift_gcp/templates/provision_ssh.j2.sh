#!/bin/bash

set -euo pipefail

if [[ -n "{{ openshift_gcp_ssh_private_key }}" ]]; then
    # Create SSH key for GCE
    if [ ! -f "{{ openshift_gcp_ssh_private_key }}" ]; then
        ssh-keygen -t rsa -f "{{ openshift_gcp_ssh_private_key }}" -C gce-provision-cloud-user -N ''
        ssh-add "{{ openshift_gcp_ssh_private_key }}" || true
    fi

    # Check if the public key is in the project metadata, and if not, add it there
    if [ -f "{{ openshift_gcp_ssh_private_key }}.pub" ]; then
        pub_file="{{ openshift_gcp_ssh_private_key }}.pub"
        pub_key=$(cut -d ' ' -f 2 < "{{ openshift_gcp_ssh_private_key }}.pub")
    else
        keyfile="${HOME}/.ssh/google_compute_engine"
        pub_file="${keyfile}.pub"
        mkdir -p "${HOME}/.ssh"
        cp "{{ openshift_gcp_ssh_private_key }}" "${keyfile}"
        chmod 0600 "${keyfile}"
        ssh-keygen -y -f "${keyfile}" >  "${pub_file}"
        pub_key=$(cut -d ' ' -f 2 <  "${pub_file}")
    fi
    key_tmp_file='/tmp/ocp-gce-keys'
    if ! gcloud --project "{{ openshift_gcp_project }}" compute project-info describe | grep -q "$pub_key"; then
        if gcloud --project "{{ openshift_gcp_project }}" compute project-info describe | grep -q ssh-rsa; then
            gcloud --project "{{ openshift_gcp_project }}" compute project-info describe | grep ssh-rsa | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/value: //' > "$key_tmp_file"
        fi
        echo -n 'cloud-user:' >> "$key_tmp_file"
        cat "${pub_file}" >> "$key_tmp_file"
        gcloud --project "{{ openshift_gcp_project }}" compute project-info add-metadata --metadata-from-file "sshKeys=${key_tmp_file}"
        rm -f "$key_tmp_file"
    fi
fi
