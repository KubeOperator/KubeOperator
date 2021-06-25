ALTER TABLE `ko`.`ko_system_registry`
    ADD COLUMN `repo_port` INT(64) NULL AFTER `architecture`;
ALTER TABLE `ko`.`ko_system_registry`
    ADD COLUMN `registry_port` INT(64) NULL AFTER `architecture`;
ALTER TABLE `ko`.`ko_system_registry`
    ADD COLUMN `registry_hosted_port` INT(64) NULL AFTER `architecture`;

update `ko`.`ko_system_registry`
set repo_port=8081;
update `ko`.`ko_system_registry`
set registry_port=8082;
update `ko`.`ko_system_registry`
set registry_hosted_port=8083;