update
    ko_cluster_manifest
set
    `is_active` = 0
where
    `name` = 'v1.22.12-ko1';

insert into
    `ko`.`ko_cluster_manifest`(
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
        `is_active`
    )
VALUES
    (
        UUID(),
        'v1.22.14-ko1',
        'v1.22.14',
        '[{\"name\":\"kubernetes\",\"version\":\"v1.22.14\"},{\"name\":\"docker\",\"version\":\"20.10.12\"},{\"name\":\"etcd\",\"version\":\"v3.5.2\"},{\"name\":\"containerd\",\"version\":\"1.6.0\"}]',
        '[{\"name\":\"calico\",\"version\":\"v3.21.4\"},{\"name\":\"flanneld\",\"version\":\"v0.15.1\"},{\"name\":\"cilium\",\"version\":\"v1.9.5\"}]',
        '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"2.4.2"},{"name":"prometheus","version":"2.34.0"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"}, {"name":"grafana","version":"8.3.1"},{"name":"logging","version":"v7.6.2"}]',
        '[{\"name\":\"external-ceph-block\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"external-cephfs\",\"version\":\"v2.1.0-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v2.5.1\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.9.0\"}, {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}, {\"name\":\"cinder\",\"version\":\"v1.20.0\"}]',
        '[{\"name\":\"coredns\",\"version\":\"1.8.0\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.6.1\"},{\"name\":\"ingress-nginx\",\"version\":\"v1.2.1\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.8.0\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR),
        0
    );

insert into
    `ko`.`ko_cluster_manifest`(
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
        `is_active`
    )
VALUES
    (
        UUID(),
        'v1.22.16-ko1',
        'v1.22.16',
        '[{\"name\":\"kubernetes\",\"version\":\"v1.22.16\"},{\"name\":\"docker\",\"version\":\"20.10.12\"},{\"name\":\"etcd\",\"version\":\"v3.5.2\"},{\"name\":\"containerd\",\"version\":\"1.6.0\"}]',
        '[{\"name\":\"calico\",\"version\":\"v3.21.4\"},{\"name\":\"flanneld\",\"version\":\"v0.15.1\"},{\"name\":\"cilium\",\"version\":\"v1.9.5\"}]',
        '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"2.4.2"},{"name":"prometheus","version":"2.34.0"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"}, {"name":"grafana","version":"8.3.1"},{"name":"logging","version":"v7.6.2"}]',
        '[{\"name\":\"external-ceph-block\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"external-cephfs\",\"version\":\"v2.1.0-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v2.5.1\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.9.0\"}, {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}, {\"name\":\"cinder\",\"version\":\"v1.20.0\"}]',
        '[{\"name\":\"coredns\",\"version\":\"1.8.0\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.6.1\"},{\"name\":\"ingress-nginx\",\"version\":\"v1.2.1\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.8.0\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR),
        1
    );

update ko_cluster_manifest set other_vars = '[{\"name\":\"coredns\",\"version\":\"1.6.7\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.2.1\"},{\"name\":\"ingress-nginx\",\"version\":\"0.33.0\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.16.9\"},{\"name\":\"helm-v3\",\"version\":\"v3.2.4\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]' where name in ("v1.18.4-ko1", "v1.18.6-ko1", "v1.18.8-ko1", "v1.18.10-ko1");

update ko_cluster_manifest set other_vars = '[{\"name\":\"coredns\",\"version\":\"1.6.7\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.2.1\"},{\"name\":\"ingress-nginx\",\"version\":\"0.33.0\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.4.1\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]' where name in ("v1.18.12-ko1", "v1.18.14-ko1", "v1.18.15-ko1", "v1.18.18-ko1", "v1.18.20-ko1");

update ko_cluster_manifest set other_vars = '[{\"name\":\"coredns\",\"version\":\"1.7.0\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.4.8\"},{\"name\":\"ingress-nginx\",\"version\":\"0.33.0\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.6.0\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]' where name in ("v1.20.4-ko1", "v1.20.6-ko1", "v1.20.8-ko1", "v1.20.10-ko1", "v1.20.14-ko1");

update ko_cluster_manifest set other_vars = '[{\"name\":\"coredns\",\"version\":\"1.8.4\"},{\"name\":\"dns-cache\",\"version\":\"1.17.0\"},{\"name\":\"traefik\",\"version\":\"v2.6.1\"},{\"name\":\"ingress-nginx\",\"version\":\"v1.2.1\"},{\"name\":\"metrics-server\",\"version\":\"v0.5.0\"},{\"name\":\"helm-v2\",\"version\":\"v2.17.0\"},{\"name\":\"helm-v3\",\"version\":\"v3.8.0\"},{\"name\":\"istio\",\"version\":\"v1.11.8\"},{\"name\":\"npd\",\"version\":\"v0.8.1\"},{\"name\":\"metallb\",\"version\":\"v0.13.7\"}]' where name in ("v1.22.6-ko1", "v1.22.8-ko1", "v1.22.10-ko1", "v1.22.12-ko1", "v1.22.14-ko1", "v1.22.16-ko1");