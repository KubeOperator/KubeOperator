
CREATE TABLE `ko_cluster_velero` (
  `id` varchar(255) NOT NULL,
  `cluster` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `bucket` varchar(255) DEFAULT NULL,
  `endpoint` varchar(255) DEFAULT NULL,
  `backup_account_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;