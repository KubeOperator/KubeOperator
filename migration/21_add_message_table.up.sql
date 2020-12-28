CREATE TABLE IF NOT EXISTS  `ko_message` (
  `id` varchar(64) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `sender` varchar(255) DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `level` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

CREATE TABLE IF NOT EXISTS  `ko_user_message` (
  `id` varchar(64) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `receive` varchar(255) DEFAULT NULL,
  `user_id` varchar(64) DEFAULT NULL,
  `message_id` varchar(64) DEFAULT NULL,
  `send_type` varchar(64) DEFAULT NULL,
  `send_status` varchar(64) DEFAULT NULL,
  `read_status` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

CREATE TABLE IF NOT EXISTS `ko_user_notification_config` (
  `id` varchar(64) NOT NULL,
  `vars` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
    `user_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

CREATE TABLE IF NOT EXISTS `ko_user_receiver` (
  `id` varchar(64) NOT NULL,
  `user_id` varchar(64) DEFAULT NULL,
  `vars` varchar(255) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

