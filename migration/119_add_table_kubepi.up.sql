CREATE TABLE IF NOT EXISTS `ko_kubepi_bind` (
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `id` varchar(255) NOT NULL,
    `source_type` varchar(64) DEFAULT NULL,
    `project` varchar(64) DEFAULT NULL,
    `cluster` varchar(64) DEFAULT NULL,
    `bind_user` varchar(64) DEFAULT NULL,
    `bind_password` varchar(64) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

insert into
    `ko`.`ko_kubepi_bind`(
        `id`,
        `source_type`,
        `project`,
        `cluster`,
        `bind_user`,
        `bind_password`,
        `created_at`,
        `updated_at`
    )
VALUES
    (
        UUID(),
        'ADMIN',
        '',
        '',
        'admin',
        'TVABAQEBAQEELTvAQm69N0AK2UwxQ4/6JHM2lUbG57A=',
        date_add(now(), interval 8 HOUR),
        date_add(now(), interval 8 HOUR)
    );