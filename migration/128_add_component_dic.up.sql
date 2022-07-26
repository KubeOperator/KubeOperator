CREATE TABLE IF NOT EXISTS `ko_component_dic` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `name` varchar(255) DEFAULT NULL,
    `type` varchar(255) DEFAULT NULL,
    `version` varchar(255) DEFAULT NULL,
    `describe` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'metrics-server', 'Metrics Server', 'v0.5.0', 'METRICS_SERVER_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'traefik', 'Ingress Controller', 'v2.2.1', 'TRAEFIK_HELPER');
INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'traefik', 'Ingress Controller', 'v2.4.8', 'TRAEFIK_HELPER');
INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'traefik', 'Ingress Controller', 'v2.6.1', 'TRAEFIK_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'ingress-nginx', 'Ingress Controller', '0.33.0', 'NGINX_HELPER');
INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'ingress-nginx', 'Ingress Controller', 'v1.1.1', 'NGINX_HELPER');
INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'ingress-nginx', 'Ingress Controller', 'v1.2.1', 'NGINX_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'gpu', 'GPU', 'v1.7.0', 'GPU_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'dns-cache', 'Dns Cache', '1.17.0', 'DNS_CACHE_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'istio', 'Istio', 'v1.11.8', 'ISTIO_HELPER');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'npd', 'Npd', 'v0.8.1', 'NPD_HELPER');