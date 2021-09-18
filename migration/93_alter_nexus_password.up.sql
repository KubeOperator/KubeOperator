
ALTER TABLE `ko`.`ko_system_registry` ADD COLUMN `nexus_password` VARCHAR(255) NULL AFTER `registry_hosted_port`;
UPDATE `ko_system_registry` SET `nexus_password`='admin123'