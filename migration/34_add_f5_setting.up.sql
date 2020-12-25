CREATE TABLE IF NOT EXISTS `ko_f5_setting` (
    `id` varchar(64)  NOT NULL,
    `cluster_id` varchar(255) DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `url` varchar(255)  DEFAULT NULL,
    `user` varchar(64)  DEFAULT NULL,
    `partition` varchar(255)  DEFAULT NULL,
    `public_ip` varchar(64)  DEFAULT NULL,
    `status` varchar(64)  DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ;