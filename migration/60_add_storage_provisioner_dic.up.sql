CREATE TABLE IF NOT EXISTS `ko_storage_provisioner_dic` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `architecture` varchar(255) DEFAULT NULL,
  `vars` mediumtext,
  PRIMARY KEY (`id`)
);

INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'external-ceph', 'v2.1.1-k8s1.11', 'amd64', '{\"rbd_provisioner_version\":\"v2.1.1-k8s1.11\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'nfs', 'v3.1.0-k8s1.11', 'all', '{\"nfs_provisioner_version\":\"v3.1.0-k8s1.11\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'vsphere', 'v1.0.3', 'amd64', '{\"vsphere_csi_version\":\"v1.0.3\", \"govc_version\":\"v0.23.0\", \"vsphere_csi_livenessprobe_version\":\"v1.1.0\", \"vsphere_csi_attacher_version\":\"v1.2.1\", \"vsphere_csi_provisioner_version\":\"v1.4.0\", \"vsphere_csi_node_driver_registrar_version\":\"v1.2.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'rook-ceph', 'v1.3.6', 'amd64', '{\"rook_ceph_version\":\"v1.3.6\", \"ceph_version\":\"v14.2.9\", \"rook_csi_ceph_version\":\"v2.1.2\", \"rook_csi_resizer_version\":\"v0.4.0\", \"rook_csi_snapshotter_version\":\"v1.2.2\", \"rook_csi_attacher_version\":\"v2.1.0\", \"rook_csi_provisioner_version\":\"v1.4.0\", \"rook_csi_node_driver_registrar_version\":\"v1.2.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_storage_provisioner_dic`(`id`, `name`, `version`, `architecture`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'oceanstor', 'v2.2.9', 'amd64', '{\"huawei_csi_driver_version\":\"2.2.9\", \"huawei_csi_attacher_version\":\"v1.2.1\"\"huawei_csi_provisioner_version\":\"v1.4.0\"\"huawei_csi_node_driver_registrar_version\":\"v1.2.0\" }', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
