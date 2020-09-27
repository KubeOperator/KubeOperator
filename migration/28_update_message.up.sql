ALTER TABLE `ko`.`ko_message` ADD COLUMN `cluster_id` varchar(64) NULL AFTER `level`;
ALTER TABLE `ko`.`ko_message` DROP COLUMN `sender`;
ALTER TABLE `ko`.`ko_message` MODIFY COLUMN `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL AFTER `title`;