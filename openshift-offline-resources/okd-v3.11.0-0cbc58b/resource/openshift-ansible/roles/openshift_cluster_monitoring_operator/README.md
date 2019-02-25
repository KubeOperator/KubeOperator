# OpenShift Cluster Monitoring Operator

This role installs the OpenShift [Cluster Monitoring Operator](https://github.com/openshift/cluster-monitoring-operator), which manages and updates the Prometheus-based monitoring stack deployed on top of OpenShift.

### **NOTE: This component is unsupported in OCP at this time.**

## Installation

To install the monitoring operator, set this variable:

```yaml
openshift_cluster_monitoring_operator_install: true
```

To uninstall, set:

```yaml
openshift_cluster_monitoring_operator_install: false
```

## Configuring Alertmanager

The Monitoring Operator comes with a [default no-op Alertmanager configuration](./defaults/main.yml). To supply a new configuration, set:

```yaml
openshift_cluster_monitoring_operator_alertmanager_config: |+
  global:
    # ...
  route:
    # ...
  receivers:
    # ...
```

The value of the variable should be a complete [Alertmanager configuration file](https://prometheus.io/docs/alerting/configuration/).

## Monitoring new components 

To integrate a new OpenShift component with monitoring, follow the [Cluster Monitoring Operator](https://github.com/openshift/cluster-monitoring-operator) docs for contributing new components.

## License

Apache License, Version 2.0
