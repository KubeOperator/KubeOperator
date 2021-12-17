ALTER TABLE
    `ko`.`ko_cluster`
ADD
    COLUMN `node_name_rule` VARCHAR(255) NULL
AFTER
    `name`;

UPDATE
    `ko`.`ko_cluster`
SET
    `node_name_rule` = "default"
WHERE
    `source` = "local"
    OR `source` = "ko-external";