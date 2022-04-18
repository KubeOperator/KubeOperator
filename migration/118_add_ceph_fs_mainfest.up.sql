update ko_cluster_manifest set storage_vars = replace(storage_vars,"external-ceph","external-ceph-rbd");
update ko_cluster_storage_provisioner set name = "external-ceph-rbd", type = "external-ceph-rbd" where type = "external-ceph";
update ko_storage_provisioner_dic set name = "external-ceph-rbd" where name = "external-ceph";

INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'external-ceph-fs', 'v2.1.0-k8s1.11', 'amd64', '{\"fs_provisioner_version\":\"v2.1.0-k8s1.11\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
