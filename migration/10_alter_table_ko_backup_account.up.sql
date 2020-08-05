ALTER TABLE `ko_backup_account`
ADD COLUMN `created_at` datetime(0) NOT NULL AFTER `credential`,
ADD COLUMN `updated_at` datetime(0) NOT NULL AFTER `created_at`;