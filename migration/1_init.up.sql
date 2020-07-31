CREATE TABLE IF NOT EXISTS `ko_cluster` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `source` varchar(255) DEFAULT NULL,
  `spec_id` varchar(255) DEFAULT NULL,
  `secret_id` varchar(255) DEFAULT NULL,
  `status_id` varchar(255) DEFAULT NULL,
  `plan_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_log` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `cluster_id` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  `status` varchar(255) DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE IF NOT EXISTS `ko_cluster_node` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `host_id` varchar(255) DEFAULT NULL,
  `cluster_id` varchar(255) DEFAULT NULL,
  `role` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_secret` (
  `id` varchar(255) NOT NULL,
  `kubeadm_token` mediumtext,
  `kubernetes_token` mediumtext,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_spec` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `provider` varchar(255) DEFAULT NULL,
  `network_type` varchar(255) DEFAULT NULL,
  `flannel_backend` varchar(255) DEFAULT NULL,
  `calico_ipv4pool_ipip` varchar(255) DEFAULT NULL,
  `runtime_type` varchar(255) DEFAULT NULL,
  `docker_storage_dir` varchar(255) DEFAULT NULL,
  `containerd_storage_dir` varchar(255) DEFAULT NULL,
  `lb_kube_apiserver_ip` varchar(255) DEFAULT NULL,
  `kube_api_server_port` int(11) DEFAULT NULL,
  `kube_router` varchar(255) DEFAULT NULL,
  `kube_pod_subnet` varchar(255) DEFAULT NULL,
  `kube_service_subnet` varchar(255) DEFAULT NULL,
  `worker_amount` int(11) DEFAULT NULL,
  `kube_max_pods` int(11) DEFAULT NULL,
  `kube_proxy_mode` varchar(255) DEFAULT NULL,
  `ingress_controller_type` varchar(255) DEFAULT NULL,
  `architectures` varchar(255) DEFAULT NULL,
  `kubernetes_audit` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_status` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `message` mediumtext,
  `phase` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_status_condition` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `cluster_status_id` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `message` text,
  `order_num` int(11) DEFAULT NULL,
  `last_probe_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_storage_provisioner` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `type` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `message` mediumtext,
  `vars` mediumtext,
  `cluster_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_tool` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `cluster_id` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `describe` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  `logo` varchar(255) DEFAULT NULL,
  `vars` mediumtext,
  `url` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_credential` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `username` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `private_key` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_demo` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_host` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `memory` int(64) DEFAULT NULL,
  `cpu_core` int(64) DEFAULT NULL,
  `os` varchar(64) DEFAULT NULL,
  `os_version` varchar(64) DEFAULT NULL,
  `gpu_num` int(64) DEFAULT NULL,
  `gpu_info` varchar(128) DEFAULT NULL,
  `ip` varchar(128) NOT NULL,
  `port` varchar(64) DEFAULT NULL,
  `credential_id` varchar(64) DEFAULT NULL,
  `status` varchar(64) DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  `zone_id` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `ip` (`ip`)
);

CREATE TABLE IF NOT EXISTS `ko_plan` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(64) DEFAULT NULL,
  `region_id` varchar(255) DEFAULT NULL,
  `deploy_template` varchar(255) DEFAULT NULL,
  `vars` mediumtext,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_plan_zones` (
  `plan_id` varchar(64) NOT NULL,
  `zone_id` varchar(64) NOT NULL,
  `id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`plan_id`,`zone_id`)
);

CREATE TABLE IF NOT EXISTS `ko_project` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(64) NOT NULL,
  `description` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_project_member` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `project_id` varchar(64) DEFAULT NULL,
  `user_id` varchar(64) DEFAULT NULL,
  `role` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_project_resource` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `resource_type` varchar(128) DEFAULT NULL,
  `resource_id` varchar(64) DEFAULT NULL,
  `project_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_region` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `datacenter` varchar(64) DEFAULT NULL,
  `provider` varchar(64) DEFAULT NULL,
  `vars` mediumtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `ko_system_setting` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `key` varchar(256) NOT NULL,
  `value` varchar(256) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`),
  UNIQUE KEY `value` (`value`)
);

CREATE TABLE IF NOT EXISTS `ko_user` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `password` varchar(256) DEFAULT NULL,
  `email` varchar(256) NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `language` varchar(64) DEFAULT NULL,
  `is_admin` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `email` (`email`)
);

CREATE TABLE IF NOT EXISTS `ko_volume` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `host_id` varchar(64) DEFAULT NULL,
  `size` varchar(64) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_zone` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `vars` mediumtext,
  `status` varchar(64) DEFAULT NULL,
  `region_id` varchar(64) DEFAULT NULL,
  `credential_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);