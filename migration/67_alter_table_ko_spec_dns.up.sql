ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `enable_dns_cache` VARCHAR(255) NULL;
ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `dns_cache_version` VARCHAR(255) NULL;

UPDATE ko_cluster_spec SET enable_dns_cache='enable', dns_cache_version='1.17.0';