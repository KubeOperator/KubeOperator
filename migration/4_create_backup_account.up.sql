CREATE TABLE IF NOT EXISTS `ko_backup_account` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) NOT NULL,
  `region` varchar(64) NOT NULL,
  `type` varchar(64) NOT NULL,
  `status` varchar(64) NOT NULL,
  `credential` mediumtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;