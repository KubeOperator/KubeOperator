CREATE TABLE IF NOT EXISTS `ko_cluster_event` (
  `id` varchar(64) NOT NULL,
  `uid` varchar(128) DEFAULT NULL,
  `message` mediumtext,
  `kind` varchar(255) DEFAULT NULL,
  `component` varchar(255) DEFAULT NULL,
  `host` varchar(255) DEFAULT NULL,
  `namespace` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `detail` mediumtext,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;