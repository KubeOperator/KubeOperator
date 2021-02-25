CREATE TABLE IF NOT EXISTS `ko_system_registry`
(
    `created_at`        date                                             DEFAULT NULL,
    `updated_at`        date                                             DEFAULT NULL,
    `id`                varchar(254) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
    `registry_hostname` varchar(255) COLLATE utf8_bin                    DEFAULT NULL,
    `registry_protocol` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
    `architecture`      varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
    PRIMARY KEY (`id`)
);
insert ignore into ko.ko_system_setting (`created_at`,`updated_at`,`id`,`key`,`value`,`tab`) values (date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR),UUID(),'arch_type','single','SYSTEM');
