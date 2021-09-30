ALTER TABLE
    `ko`.`ko_cluster_spec`
ADD
    COLUMN `master_schedule_type` VARCHAR(255) NULL
AFTER
    `provider`;

UPDATE
    `ko_cluster_spec`
SET
    `master_schedule_type` = "enable";