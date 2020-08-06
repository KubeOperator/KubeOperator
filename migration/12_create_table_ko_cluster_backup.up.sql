CREATE TABLE IF NOT  EXISTS `ko_cluster_backup_strategy` (
  `id` varchar(64) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `cron` int(64) DEFAULT NULL,
  `save_num` int(64) DEFAULT NULL,
  `backup_accoun_id` varchar(64) DEFAULT NULL,
  `status` varchar(64) DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT  EXISTS  `ko_cluster_backup_file` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  `cluster_backup_strategy_id` varchar(64) DEFAULT NULL,
  `folder` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);