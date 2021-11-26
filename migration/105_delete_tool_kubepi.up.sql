DELETE FROM
    `ko`.`ko_cluster_tool`
WHERE
    `name` = 'kubepi';

DELETE FROM
    `ko`.`ko_cluster_tool_detail`
WHERE
    `name` = 'kubepi';

INSERT INTO
    `ko_cluster_tool`(
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
    )
SELECT
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`,
    UUID() AS `id`,
    'gatekeeper' AS `name`,
    c.id,
    'v3.7.0' AS `version`,
    'OPA Gatekeeper|OPA Gatekeeper' AS `describe`,
    'Waiting' AS `status`,
    '' AS `message`,
    'gatekeeper.jpg' AS `logo`,
    '' AS `vars`,
    0 AS `frame`,
    '' AS `url`,
    'all' AS `architecture`
FROM
    `ko_cluster` c
WHERE
    c.id NOT IN (
        SELECT
            t.cluster_id
        FROM
            ko_cluster_tool t
        WHERE
            t.name = 'gatekeeper'
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
        'gatekeeper',
        'v3.7.0',
        'v3.7.0',
        'all',
        NULL,
        '{\"post_image_repo\":\"openpolicyagent/gatekeeper-crds\",\"post_image_tag\":\"v3.7.0\",\"image_repo\":\"openpolicyagent/gatekeeper\",\"crd_image_repo\":\"openpolicyagent/gatekeeper-crds\",\"image_release\":\"v3.7.0\"}',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );

UPDATE
    ko_cluster_manifest
SET
    `tool_vars` = '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.0.0"},{"name":"kubeapps","version":"v1.10.2"},{"name":"prometheus","version":"v2.18.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"v7.3.3"},{"name":"logging","version":"v7.6.2"}]'
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
    `tool_vars` = '[{"name":"gatekeeper","version":"v3.7.0"},{"name":"loki","version":"v2.1.0"},{"name":"kubeapps","version":"v2.0.1"},{"name":"prometheus","version":"v2.20.1"},{"name":"chartmuseum","version":"v0.12.0"},{"name":"registry","version":"v2.7.1"},{"name":"grafana","version":"v7.3.3"},{"name":"logging","version":"v7.6.2"}]'
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