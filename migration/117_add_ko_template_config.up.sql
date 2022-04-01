CREATE TABLE `ko_template_config` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `config` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;