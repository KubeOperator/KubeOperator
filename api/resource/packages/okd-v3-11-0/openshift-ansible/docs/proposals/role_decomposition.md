# Scaffolding for decomposing large roles

## Why?

Currently we have roles that are very large and encompass a lot of different
components. This makes for a lot of logic required within the role, can
create complex conditionals, and increases the learning curve for the role.

## How?

Creating a guide on how to approach breaking up a large role into smaller,
component based, roles. Also describe how to develop new roles, to avoid creating
large roles.

## Proposal

Create a new guide or append to the current contributing guide a process for
identifying large roles that can be split up, and how to compose smaller roles
going forward.

### Large roles

A role should be considered for decomposition if it:

1) Configures/installs more than one product.
1) Can configure multiple variations of the same product that can live
side by side.
1) Has different entry points for upgrading and installing a product

Large roles<sup>1</sup> should be responsible for:
> 1 or composing playbooks

1) Composing smaller roles to provide a full solution such as an Openshift Master
1) Ensuring that smaller roles are called in the correct order if necessary
1) Calling smaller roles with their required variables
1) Performing prerequisite tasks that small roles may depend on being in place
(openshift_logging certificate generation for example)

### Small roles

A small role should be able to:

1) Be deployed independently of other products (this is different than requiring
being installed after other base components such as OCP)
1) Be self contained and able to determine facts that it requires to complete
1) Fail fast when facts it requires are not available or are invalid
1) "Make it so" based on provided variables and anything that may be required
as part of doing such (this should include data migrations)
1) Have a minimal set of dependencies in meta/main.yml, just enough to do its job

### Example using decomposition of openshift_logging

The `openshift_logging` role was created as a port from the deployer image for
the `3.5` deliverable. It was a large role that created the service accounts,
configmaps, secrets, routes, and deployment configs/daemonset required for each
of its different components (Fluentd, Kibana, Curator, Elasticsearch).

It was possible to configure any of the components independently of one another,
up to a point. However, it was an all of nothing installation and there was a
need from customers to be able to do things like just deploy Fluentd.

Also being able to support multiple versions of configuration files would become
increasingly messy with a large role. Especially if the components had changes
at different intervals.

#### Folding of responsibility

There was a duplicate of work within the installation of three of the four logging
components where there was a possibility to deploy both an 'operations' and
'non-operations' cluster side-by-side. The first step was to collapse that
duplicate work into a single path and allow a variable to be provided to
configure such that either possibility could be created.

#### Consolidation of responsibility

The generation of OCP objects required for each component were being created in
the same task file, all Service Accounts were created at the same time, all secrets,
configmaps, etc. The only components that were not generated at the same time were
the deployment configs and the daemonset. The second step was to make the small
roles self contained and generate their own required objects.

#### Consideration for prerequisites

Currently the Aggregated Logging stack generates its own certificates as it has
some requirements that prevent it from utilizing the OCP cert generation service.
In order to make sure that all components were able to trust one another as they
did previously, until the cert generation service can be used, the certificate
generation is being handled within the top level `openshift_logging` role and
providing the location of the generated certificates to the individual roles.

#### Snippets

[openshift_logging/tasks/install_logging.yaml](https://github.com/ewolinetz/openshift-ansible/blob/logging_component_subroles/roles/openshift_logging/tasks/install_logging.yaml)
```yaml
- name: Gather OpenShift Logging Facts
  openshift_logging_facts:
    oc_bin: "{{openshift.common.client_binary}}"
    openshift_logging_namespace: "{{openshift_logging_namespace}}"

- name: Set logging project
  oc_project:
    state: present
    name: "{{ openshift_logging_namespace }}"

- name: Create logging cert directory
  file:
    path: "{{ openshift.common.config_base }}/logging"
    state: directory
    mode: 0755
  changed_when: False
  check_mode: no

- include: generate_certs.yaml
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"

## Elasticsearch
- import_role:
    name: openshift_logging_elasticsearch
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"

- import_role:
    name: openshift_logging_elasticsearch
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"
    openshift_logging_es_ops_deployment: true
  when:
  - openshift_logging_use_ops | bool


## Kibana
- import_role:
    name: openshift_logging_kibana
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"
    openshift_logging_kibana_namespace: "{{ openshift_logging_namespace }}"
    openshift_logging_kibana_master_url: "{{ openshift_logging_master_url }}"
    openshift_logging_kibana_master_public_url: "{{ openshift_logging_master_public_url }}"
    openshift_logging_kibana_image_prefix: "{{ openshift_logging_image_prefix }}"
    openshift_logging_kibana_image_version: "{{ openshift_logging_image_version }}"
    openshift_logging_kibana_replicas: "{{ openshift_logging_kibana_replica_count }}"
    openshift_logging_kibana_es_host: "{{ openshift_logging_es_host }}"
    openshift_logging_kibana_es_port: "{{ openshift_logging_es_port }}"
    openshift_logging_kibana_image_pull_secret: "{{ openshift_logging_image_pull_secret }}"

- import_role:
    name: openshift_logging_kibana
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"
    openshift_logging_kibana_ops_deployment: true
    openshift_logging_kibana_namespace: "{{ openshift_logging_namespace }}"
    openshift_logging_kibana_master_url: "{{ openshift_logging_master_url }}"
    openshift_logging_kibana_master_public_url: "{{ openshift_logging_master_public_url }}"
    openshift_logging_kibana_image_prefix: "{{ openshift_logging_image_prefix }}"
    openshift_logging_kibana_image_version: "{{ openshift_logging_image_version }}"
    openshift_logging_kibana_image_pull_secret: "{{ openshift_logging_image_pull_secret }}"
    openshift_logging_kibana_es_host: "{{ openshift_logging_es_ops_host }}"
    openshift_logging_kibana_es_port: "{{ openshift_logging_es_ops_port }}"
    openshift_logging_kibana_nodeselector: "{{ openshift_logging_kibana_ops_nodeselector }}"
    openshift_logging_kibana_memory_limit: "{{ openshift_logging_kibana_ops_memory_limit }}"
    openshift_logging_kibana_cpu_request: "{{ openshift_logging_kibana_ops_cpu_request }}"
    openshift_logging_kibana_hostname: "{{ openshift_logging_kibana_ops_hostname }}"
    openshift_logging_kibana_replicas: "{{ openshift_logging_kibana_ops_replica_count }}"
    openshift_logging_kibana_proxy_debug: "{{ openshift_logging_kibana_ops_proxy_debug }}"
    openshift_logging_kibana_proxy_memory_limit: "{{ openshift_logging_kibana_ops_proxy_memory_limit }}"
    openshift_logging_kibana_proxy_cpu_request: "{{ openshift_logging_kibana_ops_proxy_cpu_request }}"
    openshift_logging_kibana_cert: "{{ openshift_logging_kibana_ops_cert }}"
    openshift_logging_kibana_key: "{{ openshift_logging_kibana_ops_key }}"
    openshift_logging_kibana_ca: "{{ openshift_logging_kibana_ops_ca}}"
  when:
  - openshift_logging_use_ops | bool


## Curator
- import_role:
    name: openshift_logging_curator
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"
    openshift_logging_curator_namespace: "{{ openshift_logging_namespace }}"
    openshift_logging_curator_master_url: "{{ openshift_logging_master_url }}"
    openshift_logging_curator_image_prefix: "{{ openshift_logging_image_prefix }}"
    openshift_logging_curator_image_version: "{{ openshift_logging_image_version }}"
    openshift_logging_curator_image_pull_secret: "{{ openshift_logging_image_pull_secret }}"

- import_role:
    name: openshift_logging_curator
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"
    openshift_logging_curator_ops_deployment: true
    openshift_logging_curator_namespace: "{{ openshift_logging_namespace }}"
    openshift_logging_curator_master_url: "{{ openshift_logging_master_url }}"
    openshift_logging_curator_image_prefix: "{{ openshift_logging_image_prefix }}"
    openshift_logging_curator_image_version: "{{ openshift_logging_image_version }}"
    openshift_logging_curator_image_pull_secret: "{{ openshift_logging_image_pull_secret }}"
    openshift_logging_curator_memory_limit: "{{ openshift_logging_curator_ops_memory_limit }}"
    openshift_logging_curator_cpu_request: "{{ openshift_logging_curator_ops_cpu_request }}"
    openshift_logging_curator_nodeselector: "{{ openshift_logging_curator_ops_nodeselector }}"
  when:
  - openshift_logging_use_ops | bool


## Fluentd
- import_role:
    name: openshift_logging_fluentd
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"

- include: update_master_config.yaml
```

[openshift_logging_elasticsearch/meta/main.yaml](https://github.com/ewolinetz/openshift-ansible/blob/logging_component_subroles/roles/openshift_logging_elasticsearch/meta/main.yaml)
```yaml
---
galaxy_info:
  author: OpenShift Red Hat
  description: OpenShift Aggregated Logging Elasticsearch Component
  company: Red Hat, Inc.
  license: Apache License, Version 2.0
  min_ansible_version: 2.2
  platforms:
  - name: EL
    versions:
    - 7
  categories:
  - cloud
dependencies:
- role: lib_openshift
```

[openshift_logging/meta/main.yaml](https://github.com/ewolinetz/openshift-ansible/blob/logging_component_subroles/roles/openshift_logging/meta/main.yaml)
```yaml
---
galaxy_info:
  author: OpenShift Red Hat
  description: OpenShift Aggregated Logging
  company: Red Hat, Inc.
  license: Apache License, Version 2.0
  min_ansible_version: 2.2
  platforms:
  - name: EL
    versions:
    - 7
  categories:
  - cloud
dependencies:
- role: lib_openshift
- role: openshift_facts
```

[openshift_logging/tasks/install_support.yaml - old](https://github.com/openshift/openshift-ansible/blob/master/roles/openshift_logging/tasks/install_support.yaml)
```yaml
---
# This is the base configuration for installing the other components
- name: Check for logging project already exists
  command: >
    {{ openshift.common.client_binary }} --config={{ mktemp.stdout }}/admin.kubeconfig get project {{openshift_logging_namespace}} --no-headers
  register: logging_project_result
  ignore_errors: yes
  when: not ansible_check_mode
  changed_when: no

- name: "Create logging project"
  command: >
    {{ openshift.common.client_binary }} adm --config={{ mktemp.stdout }}/admin.kubeconfig new-project {{openshift_logging_namespace}}
  when: not ansible_check_mode and "not found" in logging_project_result.stderr

- name: Create logging cert directory
  file: path={{openshift.common.config_base}}/logging state=directory mode=0755
  changed_when: False
  check_mode: no

- include: generate_certs.yaml
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"

- name: Create temp directory for all our templates
  file: path={{mktemp.stdout}}/templates state=directory mode=0755
  changed_when: False
  check_mode: no

- include: generate_secrets.yaml
  vars:
    generated_certs_dir: "{{openshift.common.config_base}}/logging"

- include: generate_configmaps.yaml

- include: generate_services.yaml

- name: Generate kibana-proxy oauth client
  template: src=oauth-client.j2 dest={{mktemp.stdout}}/templates/oauth-client.yaml
  vars:
    secret: "{{oauth_secret}}"
  when: oauth_secret is defined
  check_mode: no
  changed_when: no

- include: generate_clusterroles.yaml

- include: generate_rolebindings.yaml

- include: generate_clusterrolebindings.yaml

- include: generate_serviceaccounts.yaml

- include: generate_routes.yaml
```

# Limitations

There will always be exceptions for some of these rules, however the majority of
roles should be able to fall within these guidelines.

# Additional considerations

## Playbooks including playbooks
In some circumstances it does not make sense to have a composing role but instead
a playbook would be best for orchestrating the role flow. Decisions made regarding
playbooks including playbooks will need to be taken into consideration as part of
defining this process.
Ref: (link to rteague's presentation?)

## Role dependencies
We want to make sure that our roles do not have any extra or unnecessary dependencies
in meta/main.yml without:

1. Proposing the inclusion in a team meeting or as part of the PR review and getting agreement
1. Documenting in meta/main.yml why it is there and when it was agreed to (date)

## Avoiding overly verbose roles
When we are splitting our roles up into smaller components we want to ensure we
avoid creating roles that are, for a lack of a better term, overly verbose. What
do we mean by that? If we have `openshift_control_plane` as an example, and we were to
split it up, we would have a component for `etcd`, `docker`, and possibly for
its rpms/configs. We would want to avoid creating a role that would just create
certificates as those would make sense to be contained with the rpms and configs.
Likewise, when it comes to being able to restart the master, we wouldn't have a
role where that was its sole purpose.

The same would apply for the `etcd` and `docker` roles. Anything that is required
as part of installing `etcd` such as generating certificates, installing rpms,
and upgrading data between versions should all be contained within the single
`etcd` role.

## Enforcing standards
Certain naming standards like variable names could be verified as part of a Travis
test. If we were going to also enforce that a role either has tasks or includes
(for example) then we could create tests for that as well.

## CI tests for individual roles
If we are able to correctly split up roles, it should be possible to test role
installations/upgrades like unit tests (assuming they would be able to be installed
independently of other components).
