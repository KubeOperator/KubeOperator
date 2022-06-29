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
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'metrics-server', 'Metrics Server', 'v0.5.0', 'metrics-server');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'traefik', 'Ingress Controller', 'v2.6.1', 'traefik');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'nginx', 'Ingress Controller', 'v1.1.1', 'nginx');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'gpu', 'GPU', 'v1.7.0', 'gpu');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'dns-cache', 'Dns Cache', '1.17.0', 'dns-cache');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'istio', 'Istio', 'v1.11.8', 'istio');

INSERT INTO `ko`.`ko_component_dic` (`created_at`, `updated_at`, `id`, `name`, `type`, `version`, `describe`) VALUES 
(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'npd', 'Npd', 'v0.8.1', 'npd');