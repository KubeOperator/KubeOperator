CREATE TABLE IF NOT EXISTS `ko_license` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `content` mediumtext,
  PRIMARY KEY (`id`)
) ;