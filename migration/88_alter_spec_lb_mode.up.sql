ALTER TABLE `ko_cluster_spec` ADD COLUMN `lb_mode` varchar(64) NULL AFTER `lb_kube_apiserver_ip`;
ALTER TABLE `ko_cluster_spec` ADD COLUMN `lb_kube_apiserver_port` varchar(64) NULL AFTER `lb_kube_apiserver_ip`;
UPDATE `ko_cluster_spec` SET lb_mode='internal';