import json
from uuid import UUID
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.host import Host
from storage.models import NfsStorage, CephStorage
from cloud_provider.models import Plan
from kubeops_api.models.backup_storage import BackupStorage


class JsonResourceEncoder(json.JSONEncoder):

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
                "workerSize":obj.worker_size,
                "package":obj.package.name
            }
            return cluster_data
        if isinstance(obj,Host):
            host_data = {
                "id": obj.id,
                "name": obj.name,
                "ip": obj.ip,
                "status": obj.status,
                "cpu_core": obj.cpu_core,
                "memory": obj.memory,
                "gpu": obj.gpu,
                "gpu_num": obj.gpu_num,
                "gpu_info": obj.gpu_info,
                "cluster": obj.cluster
            }
            return host_data
        if isinstance(obj,Plan):
            plan_data = {
                "id":obj.id,
                "name":obj.name,
                "deploy_template":obj.deploy_template
            }
            return plan_data
        if isinstance(obj,BackupStorage):
            backup_data = {
                "id":obj.id,
                "name":obj.name,
                "type":obj.type,
                "status":obj.status,
                "region":obj.region
            }
            return backup_data
        if isinstance(obj,NfsStorage):
            nfs_data = {
                "id": obj.id,
                "name": obj.name,
                "status":obj.status,
                "vars": obj.vars,
                "type": "NFS"
            }
            return nfs_data
        if isinstance(obj,CephStorage):
            ceph_data = {
                "id": obj.id,
                "name": obj.name,
                "vars": obj.vars,
                "type": "Ceph"
            }
            return ceph_data

        return json.JSONEncoder.default(self, obj)
