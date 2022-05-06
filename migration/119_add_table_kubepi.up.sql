CREATE TABLE IF NOT EXISTS `ko_kubepi_bind` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255)  NOT NULL,
  `source_type` varchar(64)  DEFAULT NULL,
  `source` varchar(64)  DEFAULT NULL,
  `bind_user` varchar(64)DEFAULT NULL,
  `bind_password` varchar(64)DEFAULT NULL,
  PRIMARY KEY (`id`)
);
