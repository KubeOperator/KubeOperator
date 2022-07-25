INSERT INTO `ko`.`ko_cluster_tool_detail`(
        `id`,
        `name`,
        `version`,
        `chart_version`,
        `architecture`,
        `describe`,
        `vars`,
        `created_at`,
        `updated_at`
    )
VALUES (
        UUID(),
        'prometheus',
        '2.34.0',
        '15.10.1',
        'amd64',
        NULL,
        '{\"configmap_image_name\":\"jimmidyson/configmap-reload\",\"configmap_image_tag\":\"v0.5.0\",\"metrics_image_name\":\"dyrnq/kube-state-metrics\",\"metrics_image_tag\":\"v2.4.1\",\"exporter_image_name\":\"prometheus/node-exporter\",\"exporter_image_tag\":\"v1.3.0\",\"prometheus_image_name\":\"prometheus/prometheus\",\"prometheus_image_tag\":\"v2.34.0\"}',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

UPDATE `ko_cluster_tool` t
    LEFT JOIN `ko_cluster` c ON t.cluster_id = c.id
SET t.version = '2.34.0'
WHERE c.version in ('v1.22.6-ko1', 'v1.22.8-ko1', 'v1.22.10-ko1')
    AND t.name = "prometheus"
    AND t.status = "Waiting";

UPDATE `ko_cluster_tool` t
    LEFT JOIN `ko_cluster` c ON t.cluster_id = c.id
SET t.higher_version = '2.34.0'
WHERE c.version in ('v1.22.6-ko1', 'v1.22.8-ko1', 'v1.22.10-ko1')
    AND t.name = "prometheus"
    AND t.status != "Waiting";

UPDATE `ko`.`ko_cluster_manifest`
SET `tool_vars` = '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"2.4.2"},{"name":"prometheus","version":"2.34.0"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"8.3.1"},{"name":"logging","version":"v7.6.2"}]'
WHERE `name` in (
        'v1.22.6-ko1',
        'v1.22.8-ko1',
        'v1.22.10-ko1'
    );