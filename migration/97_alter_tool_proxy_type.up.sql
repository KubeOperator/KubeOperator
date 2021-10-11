ALTER TABLE
    `ko`.`ko_cluster_tool`
ADD
    COLUMN `proxy_type` VARCHAR(255) NULL
AFTER
    `frame`;

UPDATE
    `ko_cluster_tool`
SET
    `proxy_type` = "nodeport",
    `frame` = 1
WHERE
    `name` = "prometheus"
    OR `name` = "kubepi";