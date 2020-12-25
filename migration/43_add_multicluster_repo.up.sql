CREATE TABLE IF NOT EXISTS `ko_multi_cluster_repository` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64)  NOT NULL,
  `name` varchar(256)  NOT NULL,
  `source` varchar(255)  DEFAULT NULL,
  `status` varchar(255)  DEFAULT NULL,
  `message` text ,
  `username` varchar(255)  DEFAULT NULL,
  `password` varchar(255)  DEFAULT NULL,
  `last_sync_time` datetime DEFAULT NULL,
  `last_sync_head` varchar(255)  DEFAULT NULL,
  `sync_interval` bigint DEFAULT NULL,
  `sync_status` varchar(255)  DEFAULT NULL,
  `branch` varchar(255)  DEFAULT NULL,
  `git_timeout` bigint DEFAULT NULL,
  `sync_enable` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_multi_cluster_sync_cluster_log` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64)  NOT NULL,
  `multi_cluster_sync_log_id` varchar(255)  DEFAULT NULL,
  `cluster_id` varchar(255)  DEFAULT NULL,
  `status` varchar(255)  DEFAULT NULL,
  `message` text ,
  PRIMARY KEY (`id`)
) ;


CREATE TABLE IF NOT EXISTS `ko_multi_cluster_sync_cluster_resource_log` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64)  NOT NULL,
  `resource_name` varchar(255)  DEFAULT NULL,
  `source_file` varchar(255)  DEFAULT NULL,
  `status` varchar(255)  DEFAULT NULL,
  `message` varchar(255)  DEFAULT NULL,
  `multi_cluster_sync_cluster_log_id` varchar(255)  DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

CREATE TABLE IF NOT EXISTS `ko_multi_cluster_sync_log` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64)  NOT NULL,
  `multi_cluster_repository_id` varchar(255)  DEFAULT NULL,
  `status` varchar(255)  DEFAULT NULL,
  `message` text ,
  `git_commit_id` varchar(255)  DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

CREATE TABLE IF NOT EXISTS `ko_cluster_multi_cluster_repository` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64)  NOT NULL,
  `status` varchar(255)  DEFAULT NULL,
  `message` text ,
  `multi_cluster_repository_id` varchar(255)  DEFAULT NULL,
  `cluster_id` varchar(255)  DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;