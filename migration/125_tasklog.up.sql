CREATE TABLE IF NOT EXISTS `ko_task_log` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `type` varchar(255) DEFAULT NULL,
    `phase` varchar(255) NOT NULL,
    `message` mediumtext NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_task_log_detail` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `task` varchar(255) DEFAULT NULL,
    `task_log_id` varchar(255) DEFAULT NULL,
    `last_probe_time` datetime DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext NULL,
    PRIMARY KEY (`id`)
);
