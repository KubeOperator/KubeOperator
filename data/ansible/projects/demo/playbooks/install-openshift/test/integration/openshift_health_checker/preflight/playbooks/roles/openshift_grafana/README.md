OpenShift Grafana Playbooks
===========================

OpenShift Grafana Configuration.

NOTE: Grafana is not yet supported by Red hat. This is community version of playbooks and grafana.

This role handles the configuration of Grafana dashboard with Prometheus.

Requirements
------------

* Ansible 2.2


ClusterHost Variables
--------------

For configuring new clusters, the following role variables are available.

Each host in either of the above groups must have the following variable
defined:

| Name                                         | Default value     | Description                                  |
|----------------------------------------------|-------------------|----------------------------------------------|
| openshift_grafana_namespace                  | openshift-grafana | Default grafana namespace                    |
| openshift_grafana_timeout                    | 300               | Default pod wait timeout                     |
| openshift_grafana_prometheus_namespace       | openshift-metrics | Default prometheus namespace                 |
| openshift_grafana_prometheus_serviceaccount  | prometheus        | Prometheus service account                   |
| openshift_grafana_serviceaccount_name        | grafana           | Grafana service account name                 |
| openshift_grafana_datasource_name            | prometheus        | Default datasource name                      |
| openshift_grafana_node_exporter              | false             | Do we want to deploy node exported dashboard |
| openshift_grafana_graph_granularity          | 2m                | Default dashboard granularity                |
| openshift_grafana_node_selector              | {"region":"infra"}| Default node selector                        |
| openshift_grafana_serviceaccount_annotations | empty             | Additional service account annotation list   |
| openshift_grafana_dashboards                 | (check defaults)  | Additional list of dashboards to deploy      |
| openshift_grafana_hostname                   | grafana           | Grafana route hostname                       |
| openshift_grafana_service_name               | grafana           | Grafana Service name                         |
| openshift_grafana_service_port               | 443               | Grafana service port                         |
| openshift_grafana_service_targetport         | 8443              | Grafana TargetPort to auth proxy             |
| openshift_grafana_container_port             | 3000              | Grafana container port                       |
| openshift_grafana_oauth_proxy_memory_requests| nil               | OAuthProxy memory request                    |
| openshift_grafana_oauth_proxy_cpu_requests   | nil               | OAuthProxy CPY request                       |
| openshift_grafana_oauth_proxy_memory_limit   | nil               | OAuthProxy Memory Limit                      |
| openshift_grafana_oauth_proxy_cpu_limit      | nil               | OAuthProxy CPY limit                         |
| openshift_grafana_storage_type               | emptydir          | Default grafana storage type [emptydir, pvc] |
| openshift_grafana_pvc_name                   | grafana           | Grafana Storage Claim name                   |
| openshift_grafana_pvc_access_modes           | ReadWriteOnce     | Grafana Storage Claim mode                   |
| openshift_grafana_pvc_pv_selector            | {}                | Grafana PV Selector                          |
| openshift_grafana_sc_name                    | None              | StorageClass name to use                     |

Dependencies
------------

* openshift_hosted_facts
* openshift_repos
* lib_openshift

Example Playbook
----------------

```
- name: Configure Grafana
  hosts: oo_first_master
  roles:
  - role: openshift_grafana
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Mangirdas Judeikis (mudeiki@redhat.com)
Eldad Marciano
