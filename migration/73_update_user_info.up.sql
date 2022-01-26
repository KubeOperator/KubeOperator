ALTER TABLE `ko`.`ko_user` ADD COLUMN `is_first` tinyint(1) DEFAULT '1' AFTER `type`;
ALTER TABLE `ko`.`ko_user` ADD COLUMN `err_count` int(11) DEFAULT 0 AFTER `type`;

ALTER TABLE `ko`.`ko_user` DROP INDEX `email`;