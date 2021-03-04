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
SELECT UUID(),DATE(NOW()),DATE(NOW()),
       MAX(CASE `key` WHEN 'ip' THEN value ELSE 0 END ) hostname,
       MAX(CASE `key` WHEN 'REGISTRY_PROTOCOL' THEN value ELSE 0 END ) protocol
FROM ko_system_setting;
