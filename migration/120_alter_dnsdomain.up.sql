ALTER TABLE
    `ko`.`ko_cluster_spec`
ADD
    COLUMN `kube_dns_domain` VARCHAR(255) NULL
AFTER
    `architectures`;

UPDATE
    `ko`.`ko_cluster_spec`
SET
    `kube_dns_domain` = "cluster.local";