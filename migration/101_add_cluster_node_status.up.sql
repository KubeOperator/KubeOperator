ALTER TABLE
    `ko`.`ko_cluster_status`
ADD
    COLUMN `node_cluster_id` VARCHAR(255) NULL
AFTER
    `id`;

ALTER TABLE
    `ko`.`ko_cluster_node`
ADD
    COLUMN `status_id` VARCHAR(255) NULL
AFTER
    `status`;