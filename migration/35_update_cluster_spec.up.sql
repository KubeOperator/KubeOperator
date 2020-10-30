ALTER TABLE `ko`.`ko_cluster_spec`
ADD COLUMN `network_interface` varchar(255) NULL AFTER `helm_version`;