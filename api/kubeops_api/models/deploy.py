import json
import logging

from ansible_api.models.mixins import AbstractProjectResourceModel, AbstractExecutionModel
from django.db import models

from common import models as common_models
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.setting import Setting
from kubeops_api.signals import pre_deploy_execution_start, post_deploy_execution_start

__all__ = ['DeployExecution']
logger = logging.getLogger(__name__)


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    operation = models.CharField(max_length=128, blank=False, null=False)
    progress = models.FloatField(default=0)
    current_play = models.CharField(max_length=512, null=True, default=None)
    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)

    @property
    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        cluster = Cluster.objects.get(id=self.project.id)
        hostname = Setting.objects.get(key='local_hostname')
        cluster.status = Cluster.status = Cluster.CLUSTER_STATUS_INSTALLING
        cluster.save()
        template = None
        extra_vars = {
            "cluster_name": cluster.name,
            "local_hostname": hostname.value
        }
        for temp in cluster.package.meta.get('templates', []):
            if temp['name'] == cluster.template:
                template = temp
        try:
            if len(cluster.node_set.all()) == 0:
                print("\n>>> Start Create nodes... ")
                cluster.create_resource()
                print("\n>>> End Create nodes... ")
            status_set = []
            for opt in template.get('operations', []):
                if opt['name'] == self.operation:
                    status_set = opt['status_change']
                    playbooks = []
                    cluster_playbooks = opt.get('playbooks', [])
                    playbooks.extend(cluster_playbooks)
                    total_palybook = len(playbooks)
                    current = 0
                    cluster.status = status_set['on']
                    cluster.save()
                    for playbook_name in playbooks:
                        print("\n>>> Start run {} ".format(playbook_name))
                        self.current_play = playbook_name
                        self.save()
                        playbook = self.project.playbook_set.get(name=playbook_name)
                        _result = playbook.execute(extra_vars=extra_vars)
                        result["summary"].update(_result["summary"])
                        if not _result.get('summary', {}).get('success', False):
                            cluster.status = status_set['failed']
                            cluster.save()
                            break
                        current = current + 1
                        self.progress = current / total_palybook * 100
                        self.save()
            cluster.status = status_set['succeed']
            cluster.save()
        except Exception as e:
            logger.error(e, exc_info=True)
            cluster.status = Cluster.CLUSTER_STATUS_ERROR
            cluster.save()
            result['summary'] = {'error': 'Unexpect error occur: {}'.format(e)}
        post_deploy_execution_start.send(self.__class__, execution=self, result=result)
        return result

    def to_json(self):
        dict = {'current_play': self.current_play,
                'progress': self.progress,
                'operation': self.operation,
                'state': self.state}
        return json.dumps(dict)

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)
