INSERT INTO
    `ko`.`ko_cluster_tool_detail`(
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
VALUES
    (
        UUID(),
        'prometheus',
        '2.31.1',
        '15.0.1',
        'all',
        NULL,
        '{"configmap_image_name":"jimmidyson/configmap-reload","configmap_image_tag":"v0.5.0","metrics_image_name":"dyrnq/kube-state-metrics","metrics_image_tag":"v2.2.4","exporter_image_name":"prometheus/node-exporter","exporter_image_tag":"v1.3.0","prometheus_image_name":"prometheus/prometheus","prometheus_image_tag":"v2.31.1"}',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

INSERT INTO
    `ko`.`ko_cluster_tool_detail`(
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
VALUES
    (
        UUID(),
        'kubeapps',
        '2.4.2',
        '7.6.2',
        'amd64',
        NULL,
        '',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

INSERT INTO
    `ko`.`ko_cluster_tool_detail`(
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
VALUES
    (
        UUID(),
        'grafana',
        '8.3.1',
        '6.19.0',
        'all',
        NULL,
        '{"grafana_image_name":"grafana/grafana","grafana_image_tag":"8.3.1","busybox_image_name":"kubeoperator/busybox","busybox_image_tag":"1.31.1","curl_image_name":"curlimages/curl","curl_image_tag":"7.73.0"}',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

UPDATE
    `ko`.`ko_cluster_tool_detail`
SET
    vars = '{"grafana_image_name":"grafana/grafana","grafana_image_tag":"7.3.3","busybox_image_name":"kubeoperator/busybox","busybox_image_tag":"1.28","curl_image_name":"curlimages/curl","curl_image_tag":"7.73.0"}'
WHERE
    name = 'grafana'
    AND version = 'v7.3.3';

UPDATE `ko`.`ko_cluster_tool` SET `higher_version` = "2.4.2" WHERE `name` = "kubeapps" AND `status` != "Waiting";
UPDATE `ko`.`ko_cluster_tool` SET `version` = "2.4.2" WHERE `name` = "kubeapps" AND `status` = "Waiting";
UPDATE `ko`.`ko_cluster_tool` SET `higher_version` = "2.31.1" WHERE `name` = "prometheus" AND `status` != "Waiting";
UPDATE `ko`.`ko_cluster_tool` SET `version` = "2.31.1" WHERE `name` = "prometheus" AND `status` = "Waiting";
UPDATE `ko`.`ko_cluster_tool` SET `higher_version` = "8.3.1" WHERE `name` = "grafana" AND `status` != "Waiting";
UPDATE `ko`.`ko_cluster_tool` SET `version` = "8.3.1" WHERE `name` = "grafana" AND `status` = "Waiting";

UPDATE
    `ko`.`ko_cluster_manifest`
SET
    `tool_vars` = '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.0.0"},{"name":"kubeapps","version":"2.4.2"},{"name":"prometheus","version":"2.31.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"8.3.1"},{"name":"logging","version":"v7.6.2"}]'
WHERE
    `name` in (
        'v1.18.4-ko1',
        'v1.18.6-ko1',
        'v1.18.8-ko1',
        'v1.18.10-ko1',
        'v1.18.12-ko1',
        'v1.18.14-ko1'
    );

UPDATE
    `ko`.`ko_cluster_manifest`
SET
    `tool_vars` = '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"2.4.2"},{"name":"prometheus","version":"2.31.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"8.3.1"},{"name":"logging","version":"v7.6.2"}]'
WHERE
    `name` in (
        'v1.18.15-ko1',
        'v1.18.18-ko1',
        'v1.18.20-ko1',
        'v1.20.4-ko1',
        'v1.20.6-ko1',
        'v1.20.8-ko1',
        'v1.20.10-ko1',
        'v1.20.12-ko1'
    );