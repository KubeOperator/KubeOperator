ALTER TABLE ko_cluster_spec ADD yum_operate VARCHAR(255) NULL;
UPDATE ko_cluster_spec SET yum_operate='replace';