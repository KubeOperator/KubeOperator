ALTER TABLE `ko`.`ko_system_registry` ADD COLUMN `nexus_user` VARCHAR(255) NULL AFTER `registry_hosted_port`;

UPDATE `ko`.`ko_system_registry` SET `nexus_user` = "admin";