ALTER TABLE ko_cluster_tool ADD frame int DEFAULT 0 null;
ALTER TABLE ko_cluster_tool ADD architecture varchar(255) DEFAULT 'all' null;


UPDATE ko_cluster_tool SET frame = 0;
UPDATE ko_cluster_tool SET architecture  = 'all';

UPDATE ko_cluster_tool SET frame = 1, url = '/proxy/dashboard/{cluster_name}/root' WHERE name = 'dashboard';
UPDATE ko_cluster_tool SET frame = 1, url = '/proxy/kubeapps/{cluster_name}/root' WHERE name = 'kubeapps';



UPDATE ko_cluster_tool SET architecture = 'amd64'  WHERE name = 'kubeapps';
UPDATE ko_cluster_tool SET architecture = 'amd64'  WHERE name = 'chartmuseum';
