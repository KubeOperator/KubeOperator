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
    `kube_dns_domain` varchar(255) DEFAULT NULL,
    `kubernetes_audit` varchar(255) DEFAULT NULL,
    `nodeport_address` varchar(255) DEFAULT NULL,
    `kube_service_node_port_range` varchar(255) DEFAULT NULL,
    `enable_dns_cache` varchar(255) DEFAULT NULL,
    `dns_cache_version` varchar(255) DEFAULT NULL,
    `ingress_controller_type` varchar(255) DEFAULT NULL,
    `master_schedule_type` varchar(255) DEFAULT NULL,
    `lb_mode` varchar(255) DEFAULT NULL,
    `lb_kube_apiserver_ip` varchar(255) DEFAULT NULL,
    `kube_api_server_port` int(11) DEFAULT NULL,
    `kube_router` varchar(255) DEFAULT NULL,
    `support_gpu` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `ko_cluster_spec_runtime` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `cluster_id` varchar(255) NOT NULL,
    `runtime_type` varchar(255) DEFAULT NULL,
    `docker_storage_dir` varchar(255) DEFAULT NULL,
    `containerd_storage_dir` varchar(255) DEFAULT NULL,
    `docker_subnet` varchar(255) DEFAULT NULL,
    `helm_version` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `message` varchar(255) DEFAULT NULL,
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
    `message` varchar(255) DEFAULT NULL,
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
        `kube_dns_domain`,
        `kubernetes_audit`,
        `nodeport_address`,
        `kube_service_node_port_range`,
        `enable_dns_cache`,
        `dns_cache_version`,
        `ingress_controller_type`,
        `master_schedule_type`,
        `lb_mode`,
        `lb_kube_apiserver_ip`,
        `kube_api_server_port`,
        `kube_router`,
        `support_gpu`,
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
    s.kube_dns_domain AS `kube_dns_domain`,
    s.kubernetes_audit AS `kubernetes_audit`,
    s.nodeport_address AS `nodeport_address`,
    s.kube_service_node_port_range AS `kube_service_node_port_range`,
    s.enable_dns_cache AS `enable_dns_cache`,
    s.dns_cache_version AS `dns_cache_version`,
    s.ingress_controller_type AS `ingress_controller_type`,
    s.master_schedule_type AS `master_schedule_type`,
    s.lb_mode AS `lb_mode`,
    s.lb_kube_apiserver_ip AS `lb_kube_apiserver_ip`,
    s.kube_api_server_port AS `kube_api_server_port`,
    s.kube_router AS `kube_router`,
    s.support_gpu AS `support_gpu`,
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

ALTER TABLE
    `ko`.`ko_cluster`
ADD
    COLUMN `status` VARCHAR(255) NULL
AFTER
    `source`,
ADD
    COLUMN `message` mediumtext NULL
AFTER
    `source`,
ADD
    COLUMN `provider` VARCHAR(255) NULL
AFTER
    `source`,
ADD
    COLUMN `upgrade_version` VARCHAR(255) NULL
AFTER
    `source`,
ADD
    COLUMN `version` VARCHAR(255) NULL
AFTER
    `source`,
ADD
    COLUMN `architectures` VARCHAR(255) NULL
AFTER
    `source`;


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