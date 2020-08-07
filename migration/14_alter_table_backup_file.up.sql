ALTER TABLE `ko_cluster_backup_file`
ADD COLUMN `created_at` datetime(0) NULL AFTER `folder`,
ADD COLUMN `updated_at` datetime(0) NULL AFTER `created_at`;