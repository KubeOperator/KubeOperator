ALTER TABLE `ko`.`ko_cluster_storage_provisioner` ADD `namespace` varchar(255) AFTER `name`;

UPDATE `ko`.`ko_cluster_storage_provisioner` SET `namespace`="kube-system" WHERE `type` != "rook-ceph";
UPDATE `ko`.`ko_cluster_storage_provisioner` SET `namespace`="rook-ceph" WHERE `type` = "rook-ceph";