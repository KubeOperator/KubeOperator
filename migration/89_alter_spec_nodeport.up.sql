ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `kube_service_node_port_range` VARCHAR(255) NULL AFTER `kube_proxy_mode`;
ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `nodeport_address` VARCHAR(255) NULL AFTER `kube_proxy_mode`;
UPDATE `ko_cluster_spec` SET `kube_service_node_port_range`='30000-32767'