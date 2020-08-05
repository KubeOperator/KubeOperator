ALTER TABLE `ko_backup_account`
CHANGE COLUMN `region` `bucket` varchar(255) NOT NULL AFTER `name`;