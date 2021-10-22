UPDATE
    ko_cluster_spec
    JOIN ko_cluster ON ko_cluster_spec.id = ko_cluster.spec_id
SET
    ko_cluster_spec.kube_router = (
        SELECT
            ip
        FROM
            ko_host
        WHERE
            id = (
                SELECT
                    host_id
                FROM
                    ko_cluster_node
                WHERE
                    ko_cluster_node.cluster_id = ko_cluster.id
                    AND ko_cluster_node.role = "master"
                LIMIT
                    1
            )
    );