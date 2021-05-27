INSERT INTO `ko`.`ko_cluster_tool_detail`(`id`, `name`, `version`, `chart_version`, `architecture`, `describe`, `vars`, `created_at`, `updated_at`) VALUES (
    UUID(), "dashboard", "v2.2.0", "2.2.0", "all", NULL, "{\"dashboard_image_name\":\"kubernetesui/dashboard\",\"dashboard_image_tag\":\"v2.2.0\",\"metrics_image_name\":\"kubernetesui/metrics-scraper\",\"metrics_image_tag\":\"v1.0.6\"}", date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));

UPDATE ko_cluster_tool JOIN ko_cluster ON ko_cluster_tool.cluster_id = ko_cluster.id
JOIN ko_cluster_spec ON ko_cluster_spec.id = ko_cluster.spec_id
SET ko_cluster_tool.version = 'v2.2.0' WHERE ko_cluster_spec.version LIKE 'v1.20.%' AND ko_cluster_tool.name = 'dashboard';
