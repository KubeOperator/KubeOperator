import json
import logging
from ansible_api.models.mixins import AbstractProjectResourceModel, AbstractExecutionModel
from django.db import models
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.package import Package
from kubeops_api.models.setting import Setting
from kubeops_api.signals import pre_deploy_execution_start, post_deploy_execution_start
from common import models as common_models
from kubeops_api.models.cluster_backup import ClusterBackup
from kubeops_api.storage_client import StorageClient
from kubeops_api.models.backup_storage import BackupStorage

__all__ = ['DeployExecution']
logger = logging.getLogger(__name__)


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    operation = models.CharField(max_length=128, blank=False, null=False)
    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)
    params = common_models.JsonDictTextField(default={})
    steps = common_models.JsonListTextField(default=[], null=True)

    STEP_STAUTS_PENDING = 'pending'
    STEP_STAUTS_RUNNING = 'running'
    STEP_STAUTS_SUCCESS = 'success'
    STEP_STAUTS_ERROR = 'error'

    @property
    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        cluster = self.get_cluster()
        hostname = Setting.objects.get(key='local_hostname')
        domain_suffix = Setting.objects.get(key="domain_suffix")
        extra_vars = {
            "cluster_name": cluster.name,
            "local_hostname": hostname.value,
            "domain_suffix": domain_suffix.value
        }

        extra_vars.update(cluster.configs)
        ignore_errors = False
        return_running = False

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
                result = self.on_f5_config(extra_vars)
            elif self.operation == 'upgrade':
                cluster.change_status(Cluster.CLUSTER_STATUS_UPGRADING)
                package_name = self.params.get('package', None)
                package = Package.objects.get(name=package_name)
                extra_vars.update(package.meta.get('vars'))
                result = self.on_upgrade(extra_vars)
                cluster.upgrade_package(package_name)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'scale':
                ignore_errors = True
                return_running = True
                cluster.change_status(Cluster.CLUSTER_DEPLOY_TYPE_SCALING)
                result = self.on_scaling(extra_vars)
                cluster.exit_new_node()
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'restore':
                cluster.change_status(Cluster.CLUSTER_STATUS_RESTORING)
                cluster_backup_id = self.params.get('clusterBackupId', None)
                result = self.on_restore(extra_vars, cluster_backup_id)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)

        except Exception as e:
            print('Unexpect error occur: {}'.format(e))
            if not ignore_errors:
                cluster.change_status(Cluster.CLUSTER_STATUS_ERROR)
            if return_running:
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            logger.error(e, exc_info=True)
            result['summary'] = {'error': 'Unexpect error occur: {}'.format(e)}
        post_deploy_execution_start.send(self.__class__, execution=self, result=result, ignore_errors=ignore_errors)
        return result

    def on_install(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('install')
        self.set_step_default()
        self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_RUNNING)
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            if not cluster.node_size > 0:
                try:
                    cluster.create_resource()
                    self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_SUCCESS)
                except RuntimeError as e:
                    self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_ERROR)
                    raise e
            else:
                self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_SUCCESS)
        else:
            delete = None
            for step in self.steps:
                if step['name'] == 'create-resource':
                    delete = step
            self.steps.remove(delete)
        return self.run_playbooks(extra_vars)

    def on_scaling(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('scale')
        self.set_step_default()
        self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_RUNNING)
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            try:
                num = self.params.get('num', None)
                cluster.scale_up_to(int(num))
                self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_SUCCESS)
            except RuntimeError as e:
                self.update_current_step('create-resource', DeployExecution.STEP_STAUTS_ERROR)
                raise e
        return self.run_playbooks(extra_vars)

    def on_uninstall(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('uninstall')
        self.set_step_default()
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            try:
                self.update_current_step('uninstall', DeployExecution.STEP_STAUTS_RUNNING)
                cluster.destroy_resource()
                self.update_current_step('uninstall', DeployExecution.STEP_STAUTS_SUCCESS)
            except RuntimeError as e:
                self.update_current_step('uninstall', DeployExecution.STEP_STAUTS_ERROR)
                raise e
            return {"raw": {}, "summary": {"success": True}}
        else:
            return self.run_playbooks(extra_vars)

    def on_upgrade(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('upgrade')
        self.set_step_default()
        return self.run_playbooks(extra_vars)

    def on_f5_config(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('bigip-config')
        self.set_step_default()
        extra_vars.update(cluster.meta)
        return self.run_playbooks(extra_vars)

    def on_restore(self, extra_vars, cluster_backup_id):
        cluster_backup = ClusterBackup.objects.get(id=cluster_backup_id)
        backup_storage = BackupStorage.objects.get(id=cluster_backup.backup_storage_id)
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('cluster-restore')
        client = StorageClient(backup_storage)
        backup_file_path = cluster.name + '/' + cluster_backup.name
        if client.exists(backup_file_path):
            success = client.download_file(backup_file_path,
                                           "/etc/ansible/roles/cluster-backup/files/cluster-backup.zip")
            if success:
                return self.run_playbooks(extra_vars)
            else:
                raise Exception('download file failed!')
        else:
            raise Exception('File is not exist!')

    def run_playbooks(self, extra_vars):
        result = {"raw": {}, "summary": {}}
        for step in self.steps:
            playbook_name = step.get('playbook', None)
            if playbook_name:
                playbook = self.project.playbook_set.get(name=playbook_name)
                self.update_current_step(step['name'], DeployExecution.STEP_STAUTS_RUNNING)
                _result = playbook.execute(extra_vars=extra_vars)
                result["summary"].update(_result["summary"])
                self.update_current_step(step['name'], DeployExecution.STEP_STAUTS_SUCCESS)
                if not _result.get('summary', {}).get('success', False):
                    self.update_current_step(step['name'], DeployExecution.STEP_STAUTS_ERROR)
                    raise RuntimeError("playbook: {} error!".format(step['playbook']))
        return result

    def set_step_default(self):
        for step in self.steps:
            step['status'] = DeployExecution.STEP_STAUTS_PENDING

    def get_cluster(self):
        return Cluster.objects.get(name=self.project.name)

    def update_current_step(self, name, status):
        for step in self.steps:
            if step['name'] == name:
                step['status'] = status
                self.save()

    def to_json(self):
        dict = {
            'steps': self.steps,
            'operation': self.operation,
            'state': self.state}
        return json.dumps(dict)

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)
