ALTER TABLE ko_cis_task_result ADD cluster_id VARCHAR(255) NULL;
UPDATE ko_cis_task_result r SET r.cluster_id=(SELECT t.cluster_id FROM ko_cis_task t WHERE t.id = r.cis_task_id)