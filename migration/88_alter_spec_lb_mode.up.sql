ALTER TABLE `ko_cluster_spec` ADD COLUMN `lb_mode` varchar(64) NULL AFTER `lb_kube_apiserver_ip`;
UPDATE ko_cluster_spec SET lb_mode='internal';

UPDATE ko_cluster_spec
JOIN ko_host ON ko_host.id = (SELECT ko_cluster_node.host_id FROM ko_cluster_node WHERE (ko_cluster_node.cluster_id = (SELECT id from ko_cluster WHERE spec_id = ko_cluster_spec.id))AND ko_cluster_node.role = 'master' LIMIT 1)
SET ko_cluster_spec.lb_kube_apiserver_ip = ko_host.ip
WHERE ko_cluster_spec.lb_kube_apiserver_ip IS NOT NULL;