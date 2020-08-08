ALTER TABLE `ko_user`
MODIFY COLUMN `is_admin` tinyint(1) NULL DEFAULT 0 AFTER `language`;