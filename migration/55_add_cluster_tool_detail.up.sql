CREATE TABLE IF NOT EXISTS `ko_cluster_tool_detail` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `chart_version` varchar(255) DEFAULT NULL,
  `architecture` varchar(255) DEFAULT NULL,
  `describe` varchar(255) DEFAULT NULL,
  `vars` mediumtext,
  PRIMARY KEY (`id`)
);

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'logging', 'v7.6.2', '1.0.0', 'amd64', NULL, '{\"fluentd_image_name\":\"fluentd_elasticsearch/fluentd\",\"fluentd_image_tag\":\"v2.8.0\",\"elasticsearch_image_name\":\"elasticsearch/elasticsearch\",\"elasticsearch_image_tag\":\"7.6.2\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'grafana', 'v7.3.3', '6.1.15', 'arm64', NULL, '{\"grafana_image_name\":\"kubeoperator/grafana\",\"grafana_image_tag\":\"7.3.3-arm64\",\"busybox_image_name\":\"kubeoperator/busybox\",\"busybox_image_tag\":\"1.28-arm64\",\"curl_image_name\":\"curlimages/curl\",\"curl_image_tag\":\"7.73.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'grafana', 'v7.3.3', '6.1.15', 'amd64', NULL, '{\"grafana_image_name\":\"kubeoperator/grafana\",\"grafana_image_tag\":\"7.3.3-amd64\",\"busybox_image_name\":\"kubeoperator/busybox\",\"busybox_image_tag\":\"1.28-amd64\",\"curl_image_name\":\"curlimages/curl\",\"curl_image_tag\":\"7.73.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'loki', 'v2.0.0', '2.0.0', 'arm64', NULL, '{\"loki_image_name\":\"grafana/loki\",\"loki_image_tag\":\"2.0.0-arm64\",\"promtail_image_name\":\"grafana/promtail\",\"promtail_image_tag\":\"2.0.0-arm64\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'loki', 'v2.0.0', '2.0.0', 'amd64', NULL, '{\"loki_image_name\":\"grafana/loki\",\"loki_image_tag\":\"2.0.0-amd64\",\"promtail_image_name\":\"grafana/promtail\",\"promtail_image_tag\":\"2.0.0-amd64\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'loki', 'v2.1.0', '2.3.1', 'arm64', NULL, '{\"loki_image_name\":\"grafana/loki\",\"loki_image_tag\":\"2.1.0\",\"promtail_image_name\":\"grafana/promtail\",\"promtail_image_tag\":\"2.1.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'loki', 'v2.1.0', '2.3.1', 'amd64', NULL, '{\"loki_image_name\":\"grafana/loki\",\"loki_image_tag\":\"2.1.0\",\"promtail_image_name\":\"grafana/promtail\",\"promtail_image_tag\":\"2.1.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
    
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'prometheus', 'v2.18.1', '11.5.0', 'arm64', NULL, '{\"configmap_image_name\":\"jimmidyson/configmap-reload\",\"configmap_image_tag\":\"v0.3.0\",\"metrics_image_name\":\"carlosedp/kube-state-metrics\",\"metrics_image_tag\":\"v1.9.5\",\"exporter_image_name\":\"prom/node-exporter\",\"exporter_image_tag\":\"v0.18.1\",\"prometheus_image_name\":\"prom/prometheus\",\"prometheus_image_tag\":\"v2.18.1\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'prometheus', 'v2.18.1', '11.5.0', 'amd64', NULL, '{\"configmap_image_name\":\"jimmidyson/configmap-reload\",\"configmap_image_tag\":\"v0.3.0\",\"metrics_image_name\":\"coreos/kube-state-metrics\",\"metrics_image_tag\":\"v1.9.5\",\"exporter_image_name\":\"prom/node-exporter\",\"exporter_image_tag\":\"v0.18.1\",\"prometheus_image_name\":\"prom/prometheus\",\"prometheus_image_tag\":\"v2.18.1\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'prometheus', 'v2.20.1', '11.12.1', 'arm64', NULL, '{\"configmap_image_name\":\"jimmidyson/configmap-reload\",\"configmap_image_tag\":\"v0.4.0\",\"metrics_image_name\":\"carlosedp/kube-state-metrics\",\"metrics_image_tag\":\"v1.9.5\",\"exporter_image_name\":\"prom/node-exporter\",\"exporter_image_tag\":\"v1.0.1\",\"prometheus_image_name\":\"prom/prometheus\",\"prometheus_image_tag\":\"v2.20.1\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'prometheus', 'v2.20.1', '11.12.1', 'amd64', NULL, '{\"configmap_image_name\":\"jimmidyson/configmap-reload\",\"configmap_image_tag\":\"v0.4.0\",\"metrics_image_name\":\"coreos/kube-state-metrics\",\"metrics_image_tag\":\"v1.9.6\",\"exporter_image_name\":\"prom/node-exporter\",\"exporter_image_tag\":\"v1.0.1\",\"prometheus_image_name\":\"prom/prometheus\",\"prometheus_image_tag\":\"v2.20.1\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'dashboard', 'v2.0.3', '2.2.0', 'arm64', NULL, '{\"dashboard_image_name\":\"kubernetesui/dashboard\",\"dashboard_image_tag\":\"v2.0.3\",\"metrics_image_name\":\"kubernetesui/metrics-scraper\",\"metrics_image_tag\":\"v1.0.4\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'dashboard', 'v2.0.3', '2.2.0', 'amd64', NULL, '{\"dashboard_image_name\":\"kubernetesui/dashboard\",\"dashboard_image_tag\":\"v2.0.3\",\"metrics_image_name\":\"kubernetesui/metrics-scraper\",\"metrics_image_tag\":\"v1.0.4\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'chartmuseum', 'v0.12.0', '2.13.0', 'arm64', NULL, '{\"chartmuseum_image_name\":\"kubeoperator/chartmuseum\",\"chartmuseum_image_tag\":\"v0.12.0-arm64\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'chartmuseum', 'v0.12.0', '2.13.0', 'amd64', NULL, '{\"chartmuseum_image_name\":\"chartmuseum/chartmuseum\",\"chartmuseum_image_tag\":\"v0.12.0\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'registry', 'v2.7.1', '1.9.3', 'arm64', NULL, '{\"registry_image_name\":\"kubeoperator/registry\",\"registry_image_tag\":\"2.7.1-arm64\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'registry', 'v2.7.1', '1.9.3', 'amd64', NULL, '{\"registry_image_name\":\"kubeoperator/registry\",\"registry_image_tag\":\"2.7.1-amd64\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'kubeapps', 'v1.10.2', '3.7.2', 'amd64', NULL, '{\"postgresql_image_name\":\"postgres\",\"postgresql_image_tag\":\"11-alpine\"}', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), 'kubeapps', 'v2.0.1', '5.0.1', 'amd64', NULL, '', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
