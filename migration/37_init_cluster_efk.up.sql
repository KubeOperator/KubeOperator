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
         'logging' AS `name`,
         c.id,
         'v7.6.2' AS `version`,
         '' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         'elasticsearch.png' AS `logo`,
         '' AS `vars`,
         0 AS `frame`,
         '/proxy/logging/{cluster_name}/root' AS `url`,
         s.architectures AS `architecture`
         FROM `ko_cluster` c
LEFT JOIN ko_cluster_spec s ON c.spec_id = s.id
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_tool t WHERE t.name = 'logging');