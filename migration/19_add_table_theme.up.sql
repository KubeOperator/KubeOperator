CREATE TABLE IF NOT EXISTS `ko_theme` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255)  NOT NULL,
  `system_name` varchar(255) DEFAULT NULL,
  `logo` mediumtext,
  PRIMARY KEY (`id`)
) ;

INSERT INTO ko.ko_theme (id) VALUES ('1');