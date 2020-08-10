ALTER TABLE `ko_cluster_backup_strategy`
CHANGE COLUMN `backup_accoun_id` `backup_account_id` varchar(64)  DEFAULT NULL AFTER `save_num`;