DELETE FROM `ko`.`ko_msg_subscribe` WHERE `name`='CLUSTER_UN_INSTALL';

UPDATE `ko`.`ko_msg` SET `name`='CLUSTER_DELETE' WHERE `name`='CLUSTER_UN_INSTALL';