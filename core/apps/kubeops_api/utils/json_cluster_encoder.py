import json
from uuid import UUID
from kubeops_api.models.cluster import Cluster


class JsonClusterEncoder(json.JSONEncoder):

    def default(self, obj):
        if isinstance(obj, UUID):
            return str(obj)
        if isinstance(obj, Cluster):
            cluster_data = {
                "id" : obj.id,
                "name": obj.name,
                "template": obj.template,
                "status":obj.status,
                "configs": obj.configs,
                "workerSize":obj.worker_size
            }
            return cluster_data
        return json.JSONEncoder.default(self, obj)
