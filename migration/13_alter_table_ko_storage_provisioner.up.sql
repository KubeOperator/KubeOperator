drop index name on ko_cluster_storage_provisioner;
ALTER TABLE ko_cluster_storage_provisioner ADD UNIQUE KEY(cluster_id, name);
