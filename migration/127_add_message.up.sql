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







