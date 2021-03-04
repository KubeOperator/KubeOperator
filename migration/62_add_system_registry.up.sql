CREATE TABLE IF NOT EXISTS `ko_system_registry`
(
    `created_at`        datetime     DEFAULT NULL,
    `updated_at`        datetime     DEFAULT NULL,
    `id`                varchar(64) NOT NULL,
    `hostname` varchar(255) DEFAULT NULL,
    `protocol` varchar(255) DEFAULT NULL,
    `architecture`      varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
insert into ko.ko_system_registry(
    id,
    created_at,
    updated_at,
    hostname,
    protocol
)
SELECT * FROM (
  SELECT UUID() id,DATE(NOW()) created_at,DATE(NOW()) updated_at,
         MAX(CASE `key` WHEN 'ip' THEN value ELSE '' END ) hostname,
         MAX(CASE `key` WHEN 'REGISTRY_PROTOCOL' THEN value ELSE '' END ) protocol
  FROM ko_system_setting) t
WHERE t.hostname != '';

