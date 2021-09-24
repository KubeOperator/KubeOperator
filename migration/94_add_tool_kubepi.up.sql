INSERT INTO `ko_cluster_tool`(
    `created_at`,
    `updated_at`,
    `id`,
    `name`,
    `cluster_id`,
    `version`,
    `describe`,
    `status`,
    `message`,
    `logo`,
    `vars`,
    `frame`,
    `url`,
    `architecture`
) SELECT date_add(now(), interval 8 HOUR) AS `created_at`,
         date_add(now(), interval 8 HOUR) AS `updated_at`,
         UUID() AS `id`,
         'kubepi' AS `name`,
         c.id,
         'v1.0.1' AS `version`,
         '仪表盘|Dashboard' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         'kubepi.png' AS `logo`,
         '' AS `vars`,
         1 AS `frame`,
         '/proxy/kubepi/{cluster_name}/root' AS `url`,
         'all' AS `architecture`
         FROM `ko_cluster` c
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_tool t WHERE t.name = 'kubepi');

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'kubepi', 'v1.0.1', '0.1.0', 'all', NULL, '{\"kubepi_image_name\":\"kubeoperator/kubepi-server\",\"kubepi_image_tag\":\"v1.0.1\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

UPDATE `ko`.`ko_cluster_tool` SET `describe`='应用商店|App store' WHERE `name`='kubeapps';
UPDATE `ko`.`ko_cluster_tool` SET `describe`='日志|Logs' WHERE name='logging' OR `name`='loki';
UPDATE `ko`.`ko_cluster_tool` SET `describe`='Chart 仓库|Chart warehouse' WHERE `name`='chartmuseum';
UPDATE `ko`.`ko_cluster_tool` SET `describe`='监控|Monitor' WHERE `name`='grafana' OR `name`='prometheus';
UPDATE `ko`.`ko_cluster_tool` SET `describe`='镜像仓库|Image warehouse' WHERE `name`='registry';

DELETE FROM `ko`.`ko_cluster_tool` WHERE `name`='dashboard';

DELETE FROM `ko`.`ko_cluster_tool_detail` WHERE `name`='dashboard';
