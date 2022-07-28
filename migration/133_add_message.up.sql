create TABLE `ko_msg_account` (
  `id` varchar(64) NOT NULL,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL ,
  `config` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `ko_msg_subscribe` (
  `id` varchar(64) NOT NULL,
  `name` varchar(64) NOT NULL ,
  `type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL ,
  `config` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `resource_id` varchar(64),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`)
VALUES
	(UUID(), 'CLUSTER_INSTALL', 'SYSTEM', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),'');
INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`)
VALUES
	(UUID(), 'CLUSTER_IMPORT', 'SYSTEM', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),'');
INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`)
VALUES
	(UUID(), 'CLUSTER_UN_INSTALL', 'SYSTEM', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),'');
INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`)
VALUES
	(UUID(), 'CLUSTER_DELETE', 'SYSTEM', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),'');
INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`)
VALUES
	(UUID(), 'LICENSE_EXPIRE', 'SYSTEM', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),'');

INSERT INTO `ko_msg_subscribe` (`id`, `name`, `type`, `config`, `created_at`, `updated_at`, `resource_id`) SELECT UUID(), 'CLUSTER_OPERATOR', 'CLUSTER', '{\"dingTalk\":\"DISABLE\",\"workWeiXin\":\"DISABLE\",\"local\":\"ENABLE\",\"email\":\"DISABLE\"}',  date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR),id from ko_cluster;

CREATE TABLE `ko_user_setting` (
  `id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `ko_user_setting` (`id`, `user_id`, `msg`, `created_at`, `updated_at`)
SELECT UUID(),id,'{\"dingTalk\":{\"account\":\"\",\"receive\":\"DISABLE\"},\"email\":{\"account\":\"\",\"receive\":\"DISABLE\"},\"workWeiXin\":{\"account\":\"\",\"receive\":\"DISABLE\"},\"local\":{\"account\":\"\",\"receive\":\"ENABLE\"}}',date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR) FROM ko_user;

CREATE TABLE `ko_msg_subscribe_user` (
  `subscribe_id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `id` varchar(64) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO `ko_msg_subscribe_user` (`id`,`subscribe_id`, `user_id`)
SELECT  UUID(),sub.id subscribe_id,(select user.id from ko.ko_user user where user.name='admin') user_id from ko.ko_msg_subscribe sub;

CREATE TABLE `ko_user_msg` (
  `id` varchar(64) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `receive` varchar(255) DEFAULT NULL,
  `user_id` varchar(64) DEFAULT NULL,
  `msg_id` varchar(64) DEFAULT NULL,
  `send_type` varchar(64) DEFAULT NULL,
  `send_status` varchar(64) DEFAULT NULL,
  `read_status` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `ko_msg` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `level` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `resource_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `resource_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `resource_tyoe` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
