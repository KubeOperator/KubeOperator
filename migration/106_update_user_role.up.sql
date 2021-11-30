ALTER TABLE
    `ko`.`ko_user`
ADD
    COLUMN `is_super` tinyint(1) DEFAULT '0',
AFTER
    `is_admin`;

UPDATE `ko`.`ko_user` SET `is_super`='1' WHERE `name`='admin';