CREATE TABLE IF NOT EXISTS `ko_cis_task` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255)  NOT NULL,
  `cluster_id` varchar(255)  DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `message` mediumtext,
  `status` varchar(255)DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cis_task_result` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255)  NOT NULL,
  `cis_task_id` varchar(255)  DEFAULT NULL,
  `number` varchar(255)  DEFAULT NULL,
  `desc` varchar(255)  DEFAULT NULL,
  `remediation` mediumtext ,
  `status` varchar(255)  DEFAULT NULL,
  `scored` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
);