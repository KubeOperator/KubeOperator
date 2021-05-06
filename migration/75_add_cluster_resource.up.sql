INSERT INTO ko_cluster_resource (id,resource_type,resource_id,cluster_id,created_at,updated_at)
SELECT UUID(),'HOST',id,cluster_id,date_add(now(),INTERVAL 8 HOUR),date_add(now(),INTERVAL 8 HOUR) FROM ko_host WHERE cluster_id !="";

INSERT INTO ko_cluster_resource (id,resource_type,resource_id,cluster_id,created_at,updated_at)
SELECT UUID(),'PLAN',plan_id,id,date_add(now(),INTERVAL 8 HOUR),date_add(now(),INTERVAL 8 HOUR) FROM ko_cluster WHERE plan_id !="";

INSERT INTO ko_cluster_resource (id,resource_type,resource_id,cluster_id,created_at,updated_at)
SELECT UUID(),'BACKUP_ACCOUNT',backup_account_id,cluster_id,date_add(now(),INTERVAL 8 HOUR),date_add(now(),INTERVAL 8 HOUR) FROM ko_cluster_backup_strategy;


