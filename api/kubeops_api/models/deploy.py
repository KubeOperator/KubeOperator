import json
import logging

from ansible_api.models.mixins import AbstractProjectResourceModel, AbstractExecutionModel
from django.db import models
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
        domain_suffix = Setting.objects.get(key="domain_suffix")
        extra_vars = {
            "cluster_name": cluster.name,
            "local_hostname": hostname.value,
            "domain_suffix": domain_suffix.value
        }
        ignore_errors = False
        try:
            if self.operation == "install":
                cluster.change_status(Cluster.CLUSTER_STATUS_INSTALLING)
                result = self.on_install(extra_vars)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'uninstall':
                cluster.change_status(Cluster.CLUSTER_STATUS_DELETING)
                result = self.on_uninstall(extra_vars)
                cluster.change_status(Cluster.CLUSTER_STATUS_READY)
            elif self.operation == 'bigip-config':
                ignore_errors = True
                cluster.change_status(Cluster.CLUSTER_STATUS_INSTALLING)
                result = self.on_f5_config(extra_vars)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
        except Exception as e:
            print('Unexpect error occur: {}'.format(e))
            if not ignore_errors:
                cluster.change_status(Cluster.CLUSTER_STATUS_ERROR)
            result['summary'] = {'error': 'Unexpect error occur: {}'.format(e)}
        post_deploy_execution_start.send(self.__class__, execution=self, result=result)
        return result

    def on_install(self, extra_vars):
        cluster = self.get_cluster()
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            if not cluster.node_size > 0:
                cluster.create_resource()
        playbooks = cluster.get_playbooks('install')
        return self.run_playbooks(playbooks, extra_vars)

    def on_uninstall(self, extra_vars):
        cluster = self.get_cluster()
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            self.update_progress(0)
            self.update_current_play('destroy resource')
            cluster.destroy_resource()
            self.update_progress(100)
            return {"raw": {}, "summary": {"success": True}}
        else:
            playbooks = cluster.get_playbooks('uninstall')
            return self.run_playbooks(playbooks, extra_vars)

    def on_f5_config(self, extra_vars):
        cluster = self.get_cluster()
        extra_vars.update(cluster.meta)
        playbooks = cluster.get_playbooks('bigip-config')
        return self.run_playbooks(playbooks, extra_vars)

    def run_playbooks(self, playbooks, extra_vars):
        result = {"raw": {}, "summary": {}}
        play_total = len(playbooks)
        self.update_progress(0)
        for index, playbook_name in enumerate(playbooks):
            self.update_current_play(playbook_name)
            playbook = self.project.playbook_set.get(name=playbook_name)
            _result = playbook.execute(extra_vars=extra_vars)
            result["summary"].update(_result["summary"])
            if not _result.get('summary', {}).get('success', False):
                raise Exception("playbook: {} error!".format(playbook_name))
            progress = ((index + 1) / play_total) * 100
            self.update_progress(progress)
        return result

    def get_cluster(self):
        return Cluster.objects.get(name=self.project.name)

    def update_progress(self, p):
        self.progress = p
        self.save()

    def update_current_play(self, playbook_name):
        self.current_play = playbook_name
        self.save()

    def to_json(self):
        dict = {'current_play': self.current_play,
                'progress': self.progress,
                'operation': self.operation,
                'state': self.state}
        return json.dumps(dict)

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)
