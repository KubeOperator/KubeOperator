CREATE TABLE IF NOT EXISTS `ko_system_registry`
(
    `created_at`        datetime     DEFAULT NULL,
    `updated_at`        datetime     DEFAULT NULL,
    `id`                varchar(64) NOT NULL,
    `registry_hostname` varchar(255) DEFAULT NULL,
    `registry_protocol` varchar(255) DEFAULT NULL,
    `architecture`      varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
insert ignore into ko.ko_system_setting (`created_at`, `updated_at`, `id`, `key`, `value`, `tab`)
values (date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR), UUID(), 'arch_type', 'single', 'SYSTEM');
