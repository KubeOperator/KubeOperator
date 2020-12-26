CREATE TABLE IF NOT EXISTS  `ko_ip_pool` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `subnet` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_pool` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_ip` (
  `id` varchar(64) NOT NULL,
  `address` varchar(255)  DEFAULT NULL,
  `mask` varchar(255) DEFAULT NULL,
  `gateway` varchar(255) DEFAULT NULL,
  `dns1` varchar(255) DEFAULT NULL,
  `dns2` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `ip_pool_id` varchar(64) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`address`)
);
