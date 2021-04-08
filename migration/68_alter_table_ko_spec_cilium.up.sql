ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `cilium_version` VARCHAR(255) NULL;
ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `cilium_tunnel_mode` VARCHAR(255) NULL;
ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `cilium_native_routing_cidr` VARCHAR(255) NULL;