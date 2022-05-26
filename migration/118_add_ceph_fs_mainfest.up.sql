update ko_cluster_storage_provisioner set name = "external-ceph-rbd", type = "external-ceph-rbd" where type = "external-ceph";
update ko_storage_provisioner_dic set name = "external-ceph-rbd" where name = "external-ceph";

UPDATE ko_cluster_manifest SET storage_vars='[{\"name\":\"external-ceph-rbd\",\"version\":\"v2.1.1-k8s1.11\"}, {\"name\":\"external-ceph-fs\",\"version\":\"v2.1.0-k8s1.11\"}, {\"name\":\"nfs\",\"version\":\"v3.1.0-k8s1.11\"}, {\"name\":\"vsphere\",\"version\":\"v1.0.3\"}, {\"name\":\"rook-ceph\",\"version\":\"v1.3.6\"}, {\"name\":\"oceanstor\",\"version\":\"v2.2.9\"}, {\"name\":\"cinder\",\"version\":\"v1.20.0\"}]'; 

INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'external-ceph-fs', 'v2.1.0-k8s1.11', 'amd64', '{\"fs_provisioner_version\":\"v2.1.0-k8s1.11\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
