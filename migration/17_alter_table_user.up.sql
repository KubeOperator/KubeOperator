ALTER TABLE `ko_user` ADD COLUMN `type` varchar(64) NULL AFTER `is_admin`;
UPDATE  `ko_user` SET `type` = 'LOCAL';