OpenShift Metrics with Hawkular
====================

OpenShift Metrics Installation

Requirements
------------
This role has the following dependencies:

- Java is required on the control node to generate keystores for the Java components
- httpd-tools is required on the control node to generate various passwords for the metrics components

The following variables need to be set and will be validated:

- `openshift_metrics_hawkular_hostname`: hostname used on the hawkular metrics route.

- `openshift_metrics_project`: project (i.e. namespace) where the components will be
  deployed.


Role Variables
--------------

For default values, see [`defaults/main.yaml`](defaults/main.yaml).

- `openshift_metrics_hawkular_cert:` The certificate used for re-encrypting the route
  to Hawkular metrics.  The certificate must contain the hostname used by the route.
  The default router certificate will be used if unspecified

- `openshift_metrics_hawkular_key:` The key used with the Hawkular certificate

- `openshift_metrics_hawkular_ca:` An optional certificate used to sign the Hawkular certificate.

- `openshift_metrics_hawkular_replicas:` The number of replicas for Hawkular metrics.

- `openshift_metrics_hawkular_route_annotations`: Dictionary with annotations for the Hawkular route.

- `openshift_metrics_cassandra_replicas`: The number of Cassandra nodes to deploy for the
  initial cluster.

- `openshift_metrics_cassandra_storage_type`: Use `emptydir` for ephemeral storage (for
  testing), `pv` to use persistent volumes (which need to be created before the
  installation) or `dynamic` for dynamic persistent volumes.

- `openshift_metrics_cassandra_pvc_prefix`: The name of persistent volume claims created
  for cassandra will be this with a serial number appended to the end, starting
  from 1.

- `openshift_metrics_cassandra_pvc_size`: The persistent volume claim size for each of the
  Cassandra  nodes.

- `openshift_metrics_heapster_standalone`: Deploy only heapster, without the Hawkular Metrics and
  Cassandra components.

- `openshift_metrics_heapster_allowed_users`: A comma-separated list of CN to accept.  By
  default, this is set to allow the OpenShift service proxy to connect.  If you
  override this, make sure to add `system:master-proxy` to the list in order to
  allow horizontal pod autoscaling to function properly.

- `openshift_metrics_startup_timeout`: How long in seconds we should wait until
  Hawkular Metrics and Heapster starts up before attempting a restart.

- `openshift_metrics_duration`: How many days metrics should be stored for.

- `openshift_metrics_resolution`: How often metrics should be gathered.

- `openshift_metrics_install_hawkular_agent`: Install the Hawkular OpenShift Agent (HOSA). HOSA can be used
  to collect custom metrics from your pods. This component is currently in tech-preview and is not installed by default.

## Additional variables to control resource limits
Each metrics component (hawkular, cassandra, heapster) can specify a cpu and memory limits and requests by setting
the corresponding role variable:
```
openshift_metrics_<COMPONENT>_(limits|requests)_(memory|cpu): <VALUE>
```
e.g
```
openshift_metrics_cassandra_limits_memory: 1Gi
openshift_metrics_hawkular_requests_cpu: 100
```

Dependencies
------------
openshift_facts


Example Playbook
----------------

```
- name: Configure openshift-metrics
  hosts: oo_first_master
  roles:
  - role: openshift_metrics
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Jose David Martín (j.david.nieto@gmail.com)

Image update procedure
----------------------
An upgrade of the metrics stack from older version to newer is an automated process and should be performed by calling appropriate ansible playbook and setting required ansible variables in your inventory as documented in https://docs.openshift.org/.

Following text describes manual update of the metrics images without version upgrade. To determine the current version of images being used you can:
```
oc describe pod | grep 'Image ID:'
```
This will get the repo digest that can later be compared to the inspected image details.

A way to determine when was your image last updated:
```
$ docker images
REPOSITORY                                       TAG     IMAGE ID       CREATED             SIZE
<registry>/openshift3/origin-metrics-cassandra   v3.7    f8ad8d569e27   14 hours ago        783.7 MB

$ docker inspect 9c3597aeb39f
[
    {
        . . .
        "RepoDigests": [
            "<registry>/openshift3/metrics-cassandra@sha256:d37fc0cab268625b53a92bb98d09fcc501cfca1c68e16bac6dd98446d32ba135
        ],
        . . .
        "Config": {
            . . .
            "Labels": {
                . . .
                "build-date": "2017-10-17T16:47:44.350655",
                . . .
                "release": "0.143.4.0",
                . . .
                "url": "https://access.redhat.com/containers/#/registry.access.redhat.com/openshift3/metrics-cassandra/images/v3.7.0-0.143.4.0",
                . . .
                "version": "v3.7.0"
            }
        },
        . . .
```

Pull a new image to see if registry has any newer images with the same tag:
```
$ docker pull <registry>/openshift3/origin-metrics-cassandra:v3.7
```

If there was an update, you need to run the `docker pull` on each node.

It is recommended that you now rerun the `openshift_metrics` playbook to ensure that any necessary config changes are also picked up.

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

Tue Oct 10, 2017
- Default imagePullPolicy changed from Always to IfNotPresent
