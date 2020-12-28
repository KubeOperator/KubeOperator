CREATE TABLE IF NOT EXISTS `ko_cluster_istio` (
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `cluster_id` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `describe` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `message` mediumtext,
  `vars` mediumtext,
  PRIMARY KEY (`id`)
) ;

INSERT INTO `ko_cluster_istio`(`created_at`, `updated_at`, `id`, `name`, `cluster_id`, `version`, `describe`, `status`, `message`, `vars`
) SELECT date_add(now(), interval 8 HOUR) AS `created_at`,
         date_add(now(), interval 8 HOUR) AS `updated_at`,
         UUID() AS `id`,
         'base' AS `name`,
         c.id,
         'v1.8.0' AS `version`,
         '' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         '' AS `vars`
         FROM `ko_cluster` c
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_istio t WHERE t.name = 'base');

INSERT INTO `ko_cluster_istio`(`created_at`, `updated_at`, `id`, `name`, `cluster_id`, `version`, `describe`, `status`, `message`, `vars`
) SELECT date_add(now(), interval 8 HOUR) AS `created_at`,
         date_add(now(), interval 8 HOUR) AS `updated_at`,
         UUID() AS `id`,
         'pilot' AS `name`,
         c.id,
         'v1.8.0' AS `version`,
         '' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         '' AS `vars`
         FROM `ko_cluster` c
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_istio t WHERE t.name = 'pilot');

INSERT INTO `ko_cluster_istio`(`created_at`, `updated_at`, `id`, `name`, `cluster_id`, `version`, `describe`, `status`, `message`, `vars`
) SELECT date_add(now(), interval 8 HOUR) AS `created_at`,
         date_add(now(), interval 8 HOUR) AS `updated_at`,
         UUID() AS `id`,
         'ingress' AS `name`,
         c.id,
         'v1.8.0' AS `version`,
         '' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         '' AS `vars`
         FROM `ko_cluster` c
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_istio t WHERE t.name = 'ingress');

INSERT INTO `ko_cluster_istio`(`created_at`, `updated_at`, `id`, `name`, `cluster_id`, `version`, `describe`, `status`, `message`, `vars`
) SELECT date_add(now(), interval 8 HOUR) AS `created_at`,
         date_add(now(), interval 8 HOUR) AS `updated_at`,
         UUID() AS `id`,
         'egress' AS `name`,
         c.id,
         'v1.8.0' AS `version`,
         '' AS `describe`,
         'Waiting' AS `status`,
         '' AS `message`,
         '' AS `vars`
         FROM `ko_cluster` c
WHERE c.id NOT IN (SELECT t.cluster_id FROM ko_cluster_istio t WHERE t.name = 'egress');

