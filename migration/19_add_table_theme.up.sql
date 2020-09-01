CREATE TABLE IF NOT EXISTS `ko_theme` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `system_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `logo` mediumtext COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO ko.ko_theme (created_at, updated_at,'id') VALUES ('2020-08-27 15:14:04', '2020-08-27 18:25:20','1');