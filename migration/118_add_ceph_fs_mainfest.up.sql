update ko_cluster_storage_provisioner set name = "external-ceph-rbd", type = "external-ceph-rbd" where type = "external-ceph";
update ko_storage_provisioner_dic set name = "external-ceph-rbd" where name = "external-ceph";

UPDATE ko_cluster_manifest SET storage_vars='[{\"name\":\"external-ceph-rbd\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"external-cephfs\",\"version\":\"v2.1.0-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v1.0.3\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.9.0\"}, {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}, {\"name\":\"cinder\",\"version\":\"v1.20.0\"}]';

INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'external-cephfs', 'v2.1.0-k8s1.11', 'amd64', '{\"fs_provisioner_version\":\"v2.1.0-k8s1.11\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'rook_ceph', 'v1.9.0', 'amd64', '{\"rook_ceph_version\":\"v1.9.0\", \"ceph_version\":\"v16.2.7\", \"rook_csi_ceph_version\":\"v3.6.0\", \"rook_csi_resizer_version\":\"v1.4.0\", \"rook_csi_snapshotter_version\":\"v5.0.1\", \"rook_csi_attacher_version\":\"v3.4.0\", \"rook_csi_provisioner_version\":\"v3.1.0\", \"rook_csi_node_driver_registrar_version\":\"v2.5.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
