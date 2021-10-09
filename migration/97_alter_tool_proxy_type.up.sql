ALTER TABLE
    `ko`.`ko_cluster_tool`
ADD
    COLUMN `proxy_type` VARCHAR(255) NULL
AFTER
    `frame`;

ALTER TABLE
    `ko`.`ko_cluster_tool`
ADD
    COLUMN `proxy_port` VARCHAR(255) NULL
AFTER
    `proxy_type`;

UPDATE
    `ko_cluster_tool`
SET
    `proxy_type` = "nodeport"
WHERE
    `name` = "prometheus"
    OR `name` = "kubepi";