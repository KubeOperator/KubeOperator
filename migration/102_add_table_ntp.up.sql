create TABLE IF NOT EXISTS `ko_ntp_server` (
    `id` varchar(255) NOT NULL,
    `name` varchar(64) DEFAULT NULL,
    `address` varchar(256) DEFAULT NULL,
    `status` varchar(64) DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO
    `ko`.`ko_ntp_server`(
        `id`,
        `name`,
        `address`,
        `status`,
        `created_at`,
        `updated_at`
    )
SELECT
    UUID() AS `id`,
    'ntp_server_01' AS `name`,
    c.value,
    'enable' AS `status`,
    date_add(now(), interval 8 HOUR) AS `created_at`,
    date_add(now(), interval 8 HOUR) AS `updated_at`
FROM
    `ko_system_setting` c
WHERE
    c.key = 'ntp_server';