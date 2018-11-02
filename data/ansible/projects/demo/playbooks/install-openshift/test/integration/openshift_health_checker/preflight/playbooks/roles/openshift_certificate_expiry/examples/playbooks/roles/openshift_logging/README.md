## openshift_logging Role

### Please note this role is still a work in progress

This role is used for installing the Aggregated Logging stack. It should be run against
a single host, it will create any missing certificates and API objects that the current
[logging deployer](https://github.com/openshift/origin-aggregated-logging/tree/master/deployer) does.

This role requires that the control host it is run on has Java installed as part of keystore
generation for Elasticsearch (it uses JKS) as well as openssl to sign certificates.

As part of the installation, it is recommended that you add the Fluentd node selector label
to the list of persisted [node labels](https://docs.openshift.org/latest/install_config/install/advanced_install.html#configuring-node-host-labels).

### Required vars:

- `openshift_logging_install_logging`: When `True` the `openshift_logging` role will install Aggregated Logging.

When `openshift_logging_install_logging` is set to `False` the `openshift_logging` role will uninstall Aggregated Logging.

### Optional vars:
- `openshift_logging_purge_logging`: When `openshift_logging_install_logging` is set to 'False' to trigger uninstalation and `openshift_logging_purge_logging` is set to 'True', it will completely and irreversibly remove all logging persistent data including PVC. Defaults to 'False'.
- `openshift_logging_use_ops`: If 'True', set up a second ES and Kibana cluster for infrastructure logs. Defaults to 'False'.
- `openshift_logging_master_url`: The URL for the Kubernetes master, this does not need to be public facing but should be accessible from within the cluster. Defaults to 'https://kubernetes.default.svc.{{openshift.common.dns_domain}}'.
- `openshift_logging_master_public_url`: The public facing URL for the Kubernetes master, this is used for Authentication redirection. Defaults to 'https://{{openshift.common.public_hostname}}:{{openshift_master_api_port}}'.
- `openshift_logging_namespace`: The namespace that Aggregated Logging will be installed in. Defaults to 'logging'.
- `openshift_logging_curator_default_days`: The default minimum age (in days) Curator uses for deleting log records. Defaults to '30'.
- `openshift_logging_curator_run_hour`: The hour of the day that Curator will run at. Defaults to '0'.
- `openshift_logging_curator_run_minute`: The minute of the hour that Curator will run at. Defaults to '0'.
- `openshift_logging_curator_run_timezone`: The timezone that Curator uses for figuring out its run time. Defaults to 'UTC'.
- `openshift_logging_curator_timeout`: The timeout for each Curator operation. Defaults to 300.
- `openshift_logging_curator_script_log_level`: The script log level for Curator. Defaults to 'INFO'.
- `openshift_logging_curator_log_level`: The log level for the Curator process. Defaults to 'ERROR'.
- `openshift_logging_curator_cpu_request`: The minimum amount of CPU to allocate to Curator. Default is '100m'.
- `openshift_logging_curator_memory_limit`: The amount of memory to allocate to Curator. Unset if not specified.
- `openshift_logging_curator_nodeselector`: A map of labels (e.g. {"node":"infra","region":"west"} to select the nodes where the curator pod will land.
- `openshift_logging_image_pull_secret`: The name of an existing pull secret to link to the logging service accounts

- `openshift_logging_kibana_hostname`: The Kibana hostname. Defaults to 'kibana.example.com'.
- `openshift_logging_kibana_cpu_request`: The minimum amount of CPU to allocate to Kibana or unset if not specified.
- `openshift_logging_kibana_memory_limit`: The amount of memory to allocate to Kibana or unset if not specified.
- `openshift_logging_kibana_proxy_debug`: When "True", set the Kibana Proxy log level to DEBUG. Defaults to 'false'.
- `openshift_logging_kibana_proxy_cpu_request`: The minimum amount of CPU to allocate to Kibana proxy or unset if not specified.
- `openshift_logging_kibana_proxy_memory_limit`: The amount of memory to allocate to Kibana proxy or unset if not specified.
- `openshift_logging_kibana_replica_count`: The number of replicas Kibana should be scaled up to. Defaults to 1.
- `openshift_logging_kibana_nodeselector`: A map of labels (e.g. {"node":"infra","region":"west"} to select the nodes where the pod will land.
- `openshift_logging_kibana_edge_term_policy`: Insecure Edge Termination Policy. Defaults to Redirect.
- `openshift_logging_kibana_env_vars`: A map of environment variables to add to the kibana deployment config (e.g. {"ELASTICSEARCH_REQUESTTIMEOUT":"30000"})

- `openshift_logging_fluentd_nodeselector`: The node selector that the Fluentd daemonset uses to determine where to deploy to. Defaults to '"logging-infra-fluentd": "true"'.
- `openshift_logging_fluentd_cpu_request`: The minimum amount of CPU to allocate for Fluentd collector pods. Defaults to '100m'.
- `openshift_logging_fluentd_memory_limit`: The memory limit for Fluentd pods. Defaults to '512Mi'.
- `openshift_logging_fluentd_use_journal`: *DEPRECATED - DO NOT USE* Fluentd will automatically detect whether or not Docker is using the journald log driver.
- `openshift_logging_fluentd_journal_read_from_head`: If empty, Fluentd will use its internal default, which is false.
- `openshift_logging_fluentd_hosts`: List of nodes that should be labeled for Fluentd to be deployed to. Defaults to ['--all'].
- `openshift_logging_fluentd_buffer_queue_limit`: Buffer queue limit for Fluentd. Defaults to 1024.
- `openshift_logging_fluentd_buffer_size_limit`: Buffer chunk limit for Fluentd. Defaults to 1m.
- `openshift_logging_fluentd_file_buffer_limit`: Fluentd will set the value to the file buffer limit.  Defaults to '1Gi' per destination.

- `openshift_logging_fluentd_audit_container_engine`: When `openshift_logging_fluentd_audit_container_engine` is set to `True`, the audit log of the container engine will be collected and stored in ES.
- `openshift_logging_fluentd_audit_file`: Location of audit log file. The default is `/var/log/audit/audit.log`
- `openshift_logging_fluentd_audit_pos_file`: Location of fluentd in_tail position file for the audit log file. The default is `/var/log/audit/audit.log.pos`

- `openshift_logging_es_host`: The name of the ES service Fluentd should send logs to. Defaults to 'logging-es'.
- `openshift_logging_es_port`: The port for the ES service Fluentd should sent its logs to. Defaults to '9200'.
- `openshift_logging_es_ca`: The location of the ca Fluentd uses to communicate with its openshift_logging_es_host. Defaults to '/etc/fluent/keys/ca'.
- `openshift_logging_es_client_cert`: The location of the client certificate Fluentd uses for openshift_logging_es_host. Defaults to '/etc/fluent/keys/cert'.
- `openshift_logging_es_client_key`: The location of the client key Fluentd uses for openshift_logging_es_host. Defaults to '/etc/fluent/keys/key'.

- `openshift_logging_es_cluster_size`: The number of ES cluster members. Defaults to '1'.
- `openshift_logging_es_cpu_request`: The minimum amount of CPU to allocate for an ES pod cluster member. Defaults to 1 CPU.
- `openshift_logging_es_memory_limit`: The amount of RAM that should be assigned to ES. Defaults to '8Gi'.
- `openshift_logging_es_log_appenders`: The list of rootLogger appenders for ES logs which can be: 'file', 'console'. Defaults to 'file'.
- `openshift_logging_es_pv_selector`: A key/value map added to a PVC in order to select specific PVs.  Defaults to 'None'.
- `openshift_logging_es_pvc_storage_class_name`: The name of the storage class to use for a static PVC.  Defaults to ''.
- `openshift_logging_es_pvc_dynamic`: Whether or not to add the dynamic PVC annotation for any generated PVCs. Defaults to 'False'.
- `openshift_logging_es_pvc_size`: The requested size for the ES PVCs, when not provided the role will not generate any PVCs. Defaults to '""'.
- `openshift_logging_es_pvc_prefix`: The prefix for the generated PVCs. Defaults to 'logging-es'.
- `openshift_logging_es_recover_after_time`: The amount of time ES will wait before it tries to recover. Defaults to '5m'.
- `openshift_logging_es_storage_group`: The storage group used for ES. Defaults to '65534'.
- `openshift_logging_es_nodeselector`: A map of labels (e.g. {"node":"infra","region":"west"} to select the nodes where the pod will land.
- `openshift_logging_es_number_of_shards`: The number of primary shards for every new index created in ES. Defaults to '1'.
- `openshift_logging_es_number_of_replicas`: The number of replica shards per primary shard for every new index. Defaults to '0'.

- `openshift_logging_install_eventrouter`: Coupled with `openshift_logging_install_logging`. When both are 'True', eventrouter will be installed. When both are 'False', eventrouter will be uninstalled.
Other combinations will keep the eventrouter untouched.

Detailed eventrouter configuration can be found in
- `roles/openshift_logging_eventrouter/README.md`

When `openshift_logging_use_ops` is `True`, there are some additional vars. These work the
same as above for their non-ops counterparts, but apply to the OPS cluster instance:
- `openshift_logging_es_ops_host`: logging-es-ops
- `openshift_logging_es_ops_port`: 9200
- `openshift_logging_es_ops_ca`: /etc/fluent/keys/ca
- `openshift_logging_es_ops_client_cert`: /etc/fluent/keys/cert
- `openshift_logging_es_ops_client_key`: /etc/fluent/keys/key
- `openshift_logging_es_ops_cluster_size`: 1
- `openshift_logging_es_ops_cpu_request`: The minimum amount of CPU to allocate for an ES ops pod cluster member. Defaults to 1 CPU.
- `openshift_logging_es_ops_memory_limit`: 8Gi
- `openshift_logging_es_ops_pvc_dynamic`: False
- `openshift_logging_es_ops_pvc_size`: ""
- `openshift_logging_es_ops_pvc_prefix`: logging-es-ops
- `openshift_logging_es_ops_recover_after_time`: 5m
- `openshift_logging_es_ops_storage_group`: 65534
- `openshift_logging_kibana_ops_hostname`: The Operations Kibana hostname. Defaults to 'kibana-ops.example.com'.
- `openshift_logging_kibana_ops_cpu_request`: The minimum amount of CPU to allocate to Kibana or unset if not specified.
- `openshift_logging_kibana_ops_memory_limit`: The amount of memory to allocate to Kibana or unset if not specified.
- `openshift_logging_kibana_ops_proxy_cpu_request`: The minimum amount of CPU to allocate to Kibana proxy or unset if not specified.
- `openshift_logging_kibana_ops_proxy_memory_limit`: The amount of memory to allocate to Kibana proxy or unset if not specified.
- `openshift_logging_kibana_ops_replica_count`: The number of replicas Kibana ops should be scaled up to. Defaults to 1.

Elasticsearch can be exposed for external clients outside of the cluster.
- `openshift_logging_es_allow_external`: True (default is False) - if this is
  True, Elasticsearch will be exposed as a Route
- `openshift_logging_es_hostname`: The external facing hostname to use for
  the route and the TLS server certificate (default is "es." +
  `openshift_master_default_subdomain`)
- `openshift_logging_es_cert`: The location of the certificate Elasticsearch
  uses for the external TLS server cert (default is a generated cert)
- `openshift_logging_es_key`: The location of the key Elasticsearch
  uses for the external TLS server cert (default is a generated key)
- `openshift_logging_es_ca_ext`: The location of the CA cert for the cert
  Elasticsearch uses for the external TLS server cert (default is the internal
  CA)
Elasticsearch OPS too, if using an OPS cluster:
- `openshift_logging_es_ops_allow_external`: True (default is False) - if this is
  True, Elasticsearch will be exposed as a Route
- `openshift_logging_es_ops_hostname`: The external facing hostname to use for
  the route and the TLS server certificate (default is "es-ops." +
  `openshift_master_default_subdomain`)
- `openshift_logging_es_ops_cert`: The location of the certificate Elasticsearch
  uses for the external TLS server cert (default is a generated cert)
- `openshift_logging_es_ops_key`: The location of the key Elasticsearch
  uses for the external TLS server cert (default is a generated key)
- `openshift_logging_es_ops_ca_ext`: The location of the CA cert for the cert
  Elasticsearch uses for the external TLS server cert (default is the internal
  CA)

### mux - secure_forward listener service
- `openshift_logging_use_mux`: Default `False`.  If this is `True`, a service
  called `mux` will be deployed.  This service will act as a Fluentd
  secure_forward forwarder for the node agent Fluentd daemonsets running in the
  cluster.  This can be used to reduce the number of connections to the
  OpenShift API server, by using `mux` and configuring each node Fluentd to
  send raw logs to mux and turn off the k8s metadata plugin.  This requires the
  use of `openshift_logging_mux_client_mode` (see below).
- `openshift_logging_mux_allow_external`: Default `False`.  If this is `True`,
  the `mux` service will be deployed, and it will be configured to allow
  Fluentd clients running outside of the cluster to send logs using
  secure_forward.  This allows OpenShift logging to be used as a central
  logging service for clients other than OpenShift, or other OpenShift
  clusters.
- `openshift_logging_mux_client_mode`: Values - `minimal`, `maximal`.
  Default is unset.  Setting this value will cause the Fluentd node agent to
  send logs to mux rather than directly to Elasticsearch.  The value
  `maximal` means that Fluentd will do as much processing as possible at the
  node before sending the records to mux.  This is the current recommended
  way to use mux due to current scaling issues.
  The value `minimal` means that Fluentd will do *no* processing at all, and
  send the raw logs to mux for processing.  We do not currently recommend using
  this mode, and ansible will warn you about this.
- `openshift_logging_mux_hostname`: Default is "mux." +
  `openshift_master_default_subdomain`.  This is the hostname *external*
  clients will use to connect to mux, and will be used in the TLS server cert
  subject.
- `openshift_logging_mux_port`: 24284
- `openshift_logging_mux_external_address`: The IP address that mux will listen
 on for connections from *external* clients.  Default is the default ipv4
 interface as reported by the `ansible_default_ipv4` fact.
- `openshift_logging_mux_cpu_request`: 100m
- `openshift_logging_mux_memory_limit`: 512Mi
- `openshift_logging_mux_default_namespaces`:
  openshift_logging_mux_default_namespaces is not supported.
  use openshift_logging_mux_namespaces instead.
- `openshift_logging_mux_namespaces`: Default `[]` - additional namespaces to
  create for _external_ mux clients to associate with their logs - users will
  need to set this
- `openshift_logging_mux_buffer_queue_limit`: Default `[1024]` - Buffer queue limit for Mux.
- `openshift_logging_mux_buffer_size_limit`: Default `[1m]` - Buffer chunk limit for Mux.
- `openshift_logging_mux_file_buffer_limit`: Default `[2Gi]` per destination - Mux will
  set the value to the file buffer limit.
- `openshift_logging_mux_file_buffer_storage_type`: Default `[emptydir]` - Storage
  type for the file buffer.  One of [`emptydir`, `pvc`, `hostmount`]

- `openshift_logging_mux_file_buffer_pvc_size`: The requested size for the file buffer
  PVC, when not provided the role will not generate any PVCs. Defaults to `4Gi`.
- `openshift_logging_mux_file_buffer_pvc_dynamic`: Whether or not to add the dynamic
  PVC annotation for any generated PVCs. Defaults to 'False'.
- `openshift_logging_mux_file_buffer_pvc_pv_selector`: A key/value map added to a PVC
  in order to select specific PVs.  Defaults to 'None'.
- `openshift_logging_mux_file_buffer_pvc_prefix`: The prefix for the generated PVCs.
  Defaults to 'logging-mux'.
- `openshift_logging_mux_file_buffer_storage_group`: The storage group used for Mux.
  Defaults to '65534'.

### remote syslog forwarding
- `openshift_logging_fluentd_remote_syslog`: Set `true` to enable remote syslog forwarding, defaults to `false`
- `openshift_logging_fluentd_remote_syslog_host`: Required, hostname or IP of remote syslog server
- `openshift_logging_fluentd_remote_syslog_port`: Port of remote syslog server, defaults to `514`
- `openshift_logging_fluentd_remote_syslog_severity`: Syslog severity level, defaults to `debug`
- `openshift_logging_fluentd_remote_syslog_facility`: Syslog facility, defaults to `local0`
- `openshift_logging_fluentd_remote_syslog_remove_tag_prefix`: Remove the prefix from the tag, defaults to `''` (empty)
- `openshift_logging_fluentd_remote_syslog_tag_key`: If string specified, use this field from the record to set the key field on the syslog message
- `openshift_logging_fluentd_remote_syslog_use_record`: Set `true` to use the severity and facility from the record, defaults to `false`
- `openshift_logging_fluentd_remote_syslog_payload_key`: If string is specified, use this field from the record as the payload on the syslog message

The corresponding `openshift_logging_mux_*` parameters are below.

- `openshift_logging_mux_remote_syslog`: Set `true` to enable remote syslog forwarding, defaults to `false`
- `openshift_logging_mux_remote_syslog_host`: Required, hostname or IP of remote syslog server
- `openshift_logging_mux_remote_syslog_port`: Port of remote syslog server, defaults to `514`
- `openshift_logging_mux_remote_syslog_severity`: Syslog severity level, defaults to `debug`
- `openshift_logging_mux_remote_syslog_facility`: Syslog facility, defaults to `local0`
- `openshift_logging_mux_remote_syslog_remove_tag_prefix`: Remove the prefix from the tag, defaults to `''` (empty)
- `openshift_logging_mux_remote_syslog_tag_key`: If string specified, use this field from the record to set the key field on the syslog message
- `openshift_logging_mux_remote_syslog_use_record`: Set `true` to use the severity and facility from the record, defaults to `false`
- `openshift_logging_mux_remote_syslog_payload_key`: If string is specified, use this field from the record as the payload on the syslog message

Cri-o Formatted Container Logs
------------------------------
In order to enable cri-o logs parsing, the `openshift_logging_fluentd` role
mounts `node-config.yaml` from the host to the fluentd container to this path:
```
/etc/origin/node/node-config.yaml
```

Fluentd pod on startup automatically determines from the `node-config.yaml`
whether to setup `in_tail` plugin to parse cri-o formatted logs in
`/var/log/containers/*` based on the
`kubeletArguments -> container-runtime-endpoint` value.

Image update procedure
----------------------
An upgrade of the logging stack from older version to newer is an automated process and should be performed by calling appropriate ansible playbook and setting required ansible variables in your inventory as documented in https://docs.openshift.org/.

Following text describes manual update of the logging images without version upgrade. To determine the current version of images being used you can.
```
oc describe pod | grep 'Image ID:'
```
This will get the repo digest that can later be compared to the inspected image details.

A way to determine when was your image last updated:
```
$ docker images
REPOSITORY                              TAG     IMAGE ID       CREATED             SIZE
<registry>/openshift3/logging-fluentd   v3.7    ff2e249fc45a   About an hour ago   235.2 MB

$ docker inspect ff2e249fc45a
[
    {
        . . .
        "RepoDigests": [
            "<registry>/openshift3/logging-fluentd@sha256:4346f0aa9694f32735115705ad324803b1a6ff08343c3288f7a62c3a5cb70495"
        ],
        . . .
        "Config": {
            . . .
            "Labels": {
                . . .
                "build-date": "2017-10-12T14:38:22.414827",
                . . .
                "release": "0.143.3.0",
                . . .
                "url": "https://access.redhat.com/containers/#/registry.access.redhat.com/openshift3/logging-fluentd/images/v3.7.0-0.143.3.0",
                . . .
                "version": "v3.7.0"
            }
        },
        . . .
```

Pull a new image to see if registry has any newer images with the same tag:
```
$ docker pull <registry>/openshift3/logging-fluentd:v3.7
```

If there was an update, you need to run the `docker pull` on each node.

It is recommended that you now rerun the `openshift_logging` playbook to ensure that any necessary config changes are also picked up.

To manually redeploy your pod you can do the following:
- for a DC you can do:
```
oc rollout latest <dc_name>
```

- for a RC you can scale down and scale back up
```
oc scale --replicas=0 <rc_name>

... wait for scale down

oc scale --replicas=<original_replica_count> <rc_name>
```

- for a DS you can delete the pod or unlabel and relabel your node
```
oc delete pod --selector=<ds_selector>
```

Changelog
---------
Tue Oct 26, 2017
- Make CPU request equal limit if limit is greater then request

Tue Oct 10, 2017
- Default imagePullPolicy changed from Always to IfNotPresent
