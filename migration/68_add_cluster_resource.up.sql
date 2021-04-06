create TABLE IF NOT EXISTS `ko_cluster_resource` (
  `id` varchar(64) DEFAULT NULL,
  `resource_type` varchar(64) DEFAULT NULL,
  `resource_id` varchar(64) DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
);
