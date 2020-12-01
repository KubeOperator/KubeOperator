ALTER TABLE `ko`.`ko_host` ADD COLUMN `has_gpu` TINYINT(1) NULL ;
ALTER TABLE `ko`.`ko_cluster_spec` ADD COLUMN `support_gpu` VARCHAR(255) NULL;