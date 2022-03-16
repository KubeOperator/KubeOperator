CREATE TABLE IF NOT EXISTS `ko_cluster_gpu` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `cluster_id` varchar(255) DEFAULT NULL,
  `describe` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  `vars` mediumtext,
  PRIMARY KEY (`id`)
);