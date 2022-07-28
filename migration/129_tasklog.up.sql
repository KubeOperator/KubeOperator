ALTER TABLE `ko`.`ko_cluster_node`
    ADD COLUMN `current_task_id` VARCHAR(255) NULL AFTER `status`,
    DROP COLUMN `status_id`,
    DROP COLUMN `pre_status`;

CREATE TABLE IF NOT EXISTS `ko_task_log` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `type` varchar(255) DEFAULT NULL,
    `start_time` int DEFAULT NULL,
    `end_time` int DEFAULT NULL,
    `phase` varchar(255) NOT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_task_log_detail` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `name` varchar(255) DEFAULT NULL,
    `task` varchar(255) DEFAULT NULL,
    `task_log_id` varchar(255) DEFAULT NULL,
    `cluster_id` varchar(255) DEFAULT NULL,
    `last_probe_time` int DEFAULT NULL,
    `start_time` int DEFAULT NULL,
    `end_time` int DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_task_retry_log` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `task_log_id` varchar(255) DEFAULT NULL,
    `cluster_id` varchar(255) DEFAULT NULL,
    `last_failed_time` int DEFAULT NULL,
    `restart_time` int DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);