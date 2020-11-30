CREATE TABLE IF NOT EXISTS `ko_system_log` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `operation` varchar(256) NOT NULL,
  `operation_info` varchar(256),
  PRIMARY KEY (`id`)
);