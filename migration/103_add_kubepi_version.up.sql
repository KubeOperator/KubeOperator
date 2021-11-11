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
        'kubepi',
        'v1.2.0',
        '0.1.0',
        'all',
        NULL,
        '{\"kubepi_image_name\":\"kubeoperator/kubepi-server\",\"kubepi_image_tag\":\"v1.2.0\"}',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

UPDATE
    `ko`.`ko_cluster_tool`
SET
    `higher_version` = "v1.2.0"
WHERE
    `name` = "kubepi"
    AND `status` != "Waiting";

UPDATE
    `ko`.`ko_cluster_tool`
SET
    `version` = "v1.2.0"
WHERE
    `name` = "kubepi"
    AND `status` = "Waiting";

UPDATE
    ko_cluster_manifest
SET
    `tool_vars` = '[{"name":"kubepi","version":"v1.2.0"},{"name":"loki","version":"v2.0.0"},{"name":"kubeapps","version":"v1.10.2"},{"name":"prometheus","version":"v2.18.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"v7.3.3"},{"name":"logging","version":"v7.6.2"}]'
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
    ko_cluster_manifest
SET
    `tool_vars` = '[{"name":"kubepi","version":"v1.2.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"v2.0.1"},{"name":"prometheus","version":"v2.20.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"v7.3.3"},{"name":"logging","version":"v7.6.2"}]'
WHERE
    `name` in (
        'v1.18.15-ko1',
        'v1.18.18-ko1',
        'v1.18.20-ko1',
        'v1.20.4-ko1',
        'v1.20.6-ko1',
        'v1.20.8-ko1',
        'v1.20.10-ko1'
    );