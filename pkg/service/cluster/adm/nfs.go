package adm

import "time"

func (ca *ClusterAdm) EnsureDeployNfsTaskStart(c *Cluster) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (ca *ClusterAdm) EnsureDeployNfsProvisioner(c *Cluster) error {
	time.Sleep(1 * time.Second)
	return nil
}
