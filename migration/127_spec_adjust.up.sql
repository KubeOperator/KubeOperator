CREATE TABLE IF NOT EXISTS `ko_cluster_spec_conf` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `yum_operate` varchar(255) DEFAULT NULL,
    `max_node_num` int(11) DEFAULT NULL,
    `worker_amount` int(11) DEFAULT NULL,
    `kube_max_pods` int(11) DEFAULT NULL,
    `kube_network_node_prefix` int(11) DEFAULT NULL,
    `kube_pod_subnet` varchar(255) DEFAULT NULL,
    `kube_service_subnet` varchar(255) DEFAULT NULL,
    `kube_proxy_mode` varchar(255) DEFAULT NULL,
    `cgroup_driver` varchar(255) DEFAULT NULL,
    `kube_dns_domain` varchar(255) DEFAULT NULL,
    `kubernetes_audit` varchar(255) DEFAULT NULL,
    `nodeport_address` varchar(255) DEFAULT NULL,
    `kube_service_node_port_range` varchar(255) DEFAULT NULL,
    `etcd_data_dir` varchar(255) DEFAULT NULL,
    `etcd_snapshot_count` int DEFAULT NULL,
    `etcd_compaction_retention` int DEFAULT NULL,
    `etcd_max_request` int DEFAULT NULL,
    `etcd_quota_backend` int DEFAULT NULL,
    `master_schedule_type` varchar(255) DEFAULT NULL,
    `lb_mode` varchar(255) DEFAULT NULL,
    `lb_kube_apiserver_ip` varchar(255) DEFAULT NULL,
    `kube_api_server_port` int(11) DEFAULT NULL,
    `kube_router` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_spec_runtime` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `runtime_type` varchar(255) DEFAULT NULL,
    `docker_mirror_registry` varchar(255) DEFAULT NULL,
    `docker_remote_api` varchar(255) DEFAULT NULL,
    `docker_storage_dir` varchar(255) DEFAULT NULL,
    `containerd_storage_dir` varchar(255) DEFAULT NULL,
    `docker_subnet` varchar(255) DEFAULT NULL,
    `helm_version` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);


CREATE TABLE IF NOT EXISTS `ko_cluster_spec_component` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `name` varchar(255) DEFAULT NULL,
    `type` varchar(255) DEFAULT NULL,
    `version` varchar(255) DEFAULT NULL,
    `describe` varchar(255) DEFAULT NULL,
    `vars` mediumtext,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_spec_network` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `network_type` varchar(255) DEFAULT NULL,
    `cilium_version` varchar(255) DEFAULT NULL,
    `cilium_tunnel_mode` varchar(255) DEFAULT NULL,
    `cilium_native_routing_cidr` varchar(255) DEFAULT NULL,
    `flannel_backend` varchar(255) DEFAULT NULL,
    `calico_ipv4_pool_ipip` varchar(255) DEFAULT NULL,
    `network_interface` varchar(255) DEFAULT NULL,
    `network_cidr` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` mediumtext,
    PRIMARY KEY (`id`)
);

INSERT INTO
    `ko_cluster_spec_conf`(
        `created_at`,
        `updated_at`,
        `id`,
        `cluster_id`,
        `yum_operate`,
        `max_node_num`,
        `worker_amount`,
        `kube_max_pods`,
        `kube_network_node_prefix`,
        `kube_pod_subnet`,
        `kube_service_subnet`,
        `kube_proxy_mode`,
        `cgroup_driver`,
        `kube_dns_domain`,
        `kubernetes_audit`,
        `nodeport_address`,
        `kube_service_node_port_range`,
        `master_schedule_type`,
        `lb_mode`,
        `lb_kube_apiserver_ip`,
        `kube_api_server_port`,
        `kube_router`,
        `status`,
        `message`
    )
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    c.spec_id AS `id`,
    c.id AS `cluster_id`,
    s.yum_operate AS `yum_operate`,
    s.max_node_num AS `max_node_num`,
    s.worker_amount AS `worker_amount`,
    s.kube_max_pods AS `kube_max_pods`,
    s.kube_network_node_prefix AS `kube_network_node_prefix`,
    s.kube_pod_subnet AS `kube_pod_subnet`,
    s.kube_service_subnet AS `kube_service_subnet`,
    s.kube_proxy_mode AS `kube_proxy_mode`,
    "systemd" AS `cgroup_driver`,
    s.kube_dns_domain AS `kube_dns_domain`,
    s.kubernetes_audit AS `kubernetes_audit`,
    s.nodeport_address AS `nodeport_address`,
    s.kube_service_node_port_range AS `kube_service_node_port_range`,
    s.master_schedule_type AS `master_schedule_type`,
    s.lb_mode AS `lb_mode`,
    s.lb_kube_apiserver_ip AS `lb_kube_apiserver_ip`,
    s.kube_api_server_port AS `kube_api_server_port`,
    s.kube_router AS `kube_router`,
    "Running" AS `status`,
    "" AS `message`
FROM
    `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id;

INSERT INTO
    `ko_cluster_spec_runtime`(
        `created_at`,
        `updated_at`,
        `id`,
        `cluster_id`,
        `runtime_type`,
        `docker_storage_dir`,
        `containerd_storage_dir`,
        `docker_subnet`,
        `helm_version`,
        `status`,
        `message`
    )
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    c.spec_id AS `id`,
    c.id AS `cluster_id`,
    s.runtime_type AS `runtime_type`,
    s.docker_storage_dir AS `docker_storage_dir`,
    s.containerd_storage_dir AS `containerd_storage_dir`,
    s.docker_subnet AS `docker_subnet`,
    s.helm_version AS `helm_version`,
    "Running" AS `status`,
    "" AS `message`
FROM
    `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id;

INSERT INTO
    `ko_cluster_spec_network`(
        `created_at`,
        `updated_at`,
        `id`,
        `cluster_id`,
        `network_type`,
        `cilium_version`,
        `cilium_tunnel_mode`,
        `cilium_native_routing_cidr`,
        `flannel_backend`,
        `calico_ipv4_pool_ipip`,
        `network_interface`,
        `network_cidr`,
        `status`,
        `message`
    )
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    c.spec_id AS `id`,
    c.id AS `cluster_id`,
    s.network_type AS `network_type`,
    s.cilium_version AS `cilium_version`,
    s.cilium_tunnel_mode AS `cilium_tunnel_mode`,
    s.cilium_native_routing_cidr AS `cilium_native_routing_cidr`,
    s.flannel_backend AS `flannel_backend`,
    s.calico_ipv4pool_ipip AS `calico_ipv4_pool_ipip`,
    s.network_interface AS `network_interface`,
    s.network_cidr AS `network_cidr`,
    "Running" AS `status`,
    "" AS `message`
FROM
    `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id;

ALTER TABLE `ko`.`ko_cluster`
ADD COLUMN `status` VARCHAR(255) NULL AFTER `source`,
ADD COLUMN `current_task_id` VARCHAR(255) NULL AFTER `source`,
ADD COLUMN `message` mediumtext NULL AFTER `source`,
ADD COLUMN `provider` VARCHAR(255) NULL AFTER `source`,
ADD COLUMN `upgrade_version` VARCHAR(255) NULL AFTER `source`,
ADD COLUMN `version` VARCHAR(255) NULL AFTER `source`,
ADD COLUMN `architectures` VARCHAR(255) NULL AFTER `source`;


UPDATE
    `ko_cluster` c
    JOIN `ko_cluster_spec` s ON s.id = c.spec_id
    JOIN `ko_cluster_status` x ON x.id = c.status_id
SET
    c.provider = s.provider,
    c.upgrade_version = s.upgrade_version,
    c.version = s.version,
    c.architectures = s.architectures,
    c.status = x.phase,
    c.message = x.message;


INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "metrics-server" AS `name`,
    "Metrics Server" AS `type`,
    "v0.5.0" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c WHERE c.source = 'local' OR c.source = 'ko-external';


INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "dns-cache" AS `name`,
    "Dns Cache" AS `type`,
    "1.17.0" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.enable_dns_cache = 'enable';


INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "gpu" AS `name`,
    "GPU" AS `type`,
    "v1.7.0" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.support_gpu = 'enable';


INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "traefik" AS `name`,
    "Ingress Controller" AS `type`,
    "v2.2.1" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.ingress_controller_type = 'traefik' AND c.version in ("v1.18.4-ko1", "v1.18.6-ko1", "v1.18.8-ko1", "v1.18.10-ko1", "v1.18.12-ko1", "v1.18.14-ko1", "v1.18.15-ko1", "v1.18.18-ko1", "v1.18.20-ko1");

INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "traefik" AS `name`,
    "Ingress Controller" AS `type`,
    "v2.4.8" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.ingress_controller_type = 'traefik' AND c.version in ("v1.20.4-ko1", "v1.20.6-ko1", "v1.20.8-ko1", "v1.20.10-ko1", "v1.20.14-ko1");

INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "traefik" AS `name`,
    "Ingress Controller" AS `type`,
    "v2.6.1" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.ingress_controller_type = 'traefik' AND c.version in ("v1.22.6-ko1", "v1.22.8-ko1", "v1.22.10-ko1", "v1.22.12-ko1");

INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "ingress-nginx" AS `name`,
    "Ingress Controller" AS `type`,
    "v1.2.1" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.ingress_controller_type = 'nginx' AND c.version in ("v1.22.6-ko1", "v1.22.8-ko1", "v1.22.10-ko1");

INSERT INTO
    `ko_cluster_spec_component`(`created_at`, `updated_at`, `id`, `cluster_id`, `name`, `type`, `version`, `status`, `message`)
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    c.id AS `cluster_id`,
    "ingress-nginx" AS `name`,
    "Ingress Controller" AS `type`,
    "0.33.0" AS `version`,
    "enable" AS `status`,
    "" AS `message`
FROM `ko_cluster` c
    LEFT JOIN `ko_cluster_spec` s ON c.spec_id = s.id
    WHERE s.ingress_controller_type = 'nginx' AND c.version not in ("v1.22.6-ko1", "v1.22.8-ko1", "v1.22.10-ko1");