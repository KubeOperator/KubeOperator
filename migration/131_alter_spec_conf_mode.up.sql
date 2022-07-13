ALTER TABLE `ko`.`ko_cluster_spec_conf` ADD `authentication_mode` varchar(255) AFTER `kube_router`;
UPDATE `ko`.`ko_cluster_spec_conf` SET `authentication_mode`="bearer";

ALTER TABLE `ko`.`ko_cluster_secret` 
    ADD `config_content` mediumtext AFTER `kubernetes_token`,
    ADD `key_data_str` mediumtext AFTER `kubernetes_token`,
    ADD `cert_data_str` mediumtext AFTER `kubernetes_token`;
