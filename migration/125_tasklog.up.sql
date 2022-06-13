CREATE TABLE IF NOT EXISTS `ko_task_log` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `type` varchar(255) DEFAULT NULL,
    `start_time` datetime DEFAULT NULL,
    `end_time` datetime DEFAULT NULL,
    `phase` varchar(255) NOT NULL,
    `pre_phase`varchar(255) NOT NULL,
    `message` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_task_log_detail` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `task` varchar(255) DEFAULT NULL,
    `task_id` varchar(255) DEFAULT NULL,
    `start_time` datetime DEFAULT NULL,
    `end_time` datetime DEFAULT NULL,
    `last_probe_time` datetime DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
