update ko_cluster_manifest set is_active=0 where is_active='1';

insert into `ko`.`ko_cluster_manifest`(
    `id`,
    `name`,
    `version`,
    `core_vars`,
    `network_vars`,
    `tool_vars`,
    `storage_vars`,
    `other_vars`,
    `created_at`,
    `updated_at`,
    `is_active`)
VALUES (
           UUID(),
           'v1.18.20-ko1',
           'v1.18.20',
           '[{\"name\":\"kubernetes\",\"version\":\"v1.18.20\"},{\"name\":\"docker\",\"version\":\"19.03.15\"},{\"name\":\"etcd\",\"version\":\"v3.4.14\"},{\"name\":\"containerd\",\"version\":\"1.4.3\"}]',
           '[{\"name\":\"calico\",\"version\":\"v3.18.4\"},{\"name\":\"flanneld\",\"version\":\"v0.13.0\"},{\"name\":\"cilium\",\"version\":\"v1.9.5\"}]',
           '[{\"name\":\"dashboard\",\"version\":\"v2.2.0\"},{\"name\":\"loki\",\"version\":\"v2.1.0\"},{\"name\":\"kubeapps\",\"version\":\"v2.0.1\"},{\"name\":\"prometheus\",\"version\":\"v2.20.1\"},{\"name\":\"chartmuseum\",\"version\":\"v0.12.0\"},{\"name\":\"registry\",\"version\":\"v2.7.1\"},{\"name\":\"grafana\",\"version\":\"v7.3.3\"},{\"name\":\"logging\",\"version\":\"v7.6.2\"}]',
           '[{\"name\":\"external-ceph\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v1.0.3\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.3.6\"} , {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}]',
           '[{\"name\":\"coredns\",\"version\":\"1.7.0\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.4.8\"},{\"name\":\"ingress-nginx\",\"version\":\"0.33.0\"},{\"name\":\"metrics-server\",\"version\":\"v0.3.6\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.6.0\"}]',
           date_add(now(), interval 8 HOUR),
           date_add(now(), interval 8 HOUR),
           1);

insert into `ko`.`ko_cluster_manifest`(
    `id`,
    `name`,
    `version`,
    `core_vars`,
    `network_vars`,
    `tool_vars`,
    `storage_vars`,
    `other_vars`,
    `created_at`,
    `updated_at`,
    `is_active`)
VALUES (
           UUID(),
           'v1.20.8-ko1',
           'v1.20.8',
           '[{\"name\":\"kubernetes\",\"version\":\"v1.20.8\"},{\"name\":\"docker\",\"version\":\"19.03.15\"},{\"name\":\"etcd\",\"version\":\"v3.4.14\"},{\"name\":\"containerd\",\"version\":\"1.4.3\"}]',
           '[{\"name\":\"calico\",\"version\":\"v3.18.4\"},{\"name\":\"flanneld\",\"version\":\"v0.13.0\"},{\"name\":\"cilium\",\"version\":\"v1.9.5\"}]',
           '[{\"name\":\"dashboard\",\"version\":\"v2.2.0\"},{\"name\":\"loki\",\"version\":\"v2.1.0\"},{\"name\":\"kubeapps\",\"version\":\"v2.0.1\"},{\"name\":\"prometheus\",\"version\":\"v2.20.1\"},{\"name\":\"chartmuseum\",\"version\":\"v0.12.0\"},{\"name\":\"registry\",\"version\":\"v2.7.1\"},{\"name\":\"grafana\",\"version\":\"v7.3.3\"},{\"name\":\"logging\",\"version\":\"v7.6.2\"}]',
           '[{\"name\":\"external-ceph\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v1.0.3\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.3.6\"} , {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}]',
           '[{\"name\":\"coredns\",\"version\":\"1.7.0\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.4.8\"},{\"name\":\"ingress-nginx\",\"version\":\"0.33.0\"},{\"name\":\"metrics-server\",\"version\":\"v0.3.6\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.6.0\"}]',
           date_add(now(), interval 8 HOUR),
           date_add(now(), interval 8 HOUR),
           1);