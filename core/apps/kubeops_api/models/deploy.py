import json
import logging
from ansible_api.models.mixins import AbstractProjectResourceModel, AbstractExecutionModel
from django.db import models
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.package import Package
from kubeops_api.models.role import Role
from kubeops_api.models.setting import Setting
from kubeops_api.signals import pre_deploy_execution_start, post_deploy_execution_start
from common import models as common_models
from kubeops_api.models.cluster_backup import ClusterBackup
from kubeops_api.storage_client import StorageClient
from kubeops_api.models.backup_storage import BackupStorage
import kubeops_api.cluster_backup_utils
import kubeops_api.cluster_monitor
from django.utils import timezone
from message_center.message_client import MessageClient

__all__ = ['DeployExecution']
logger = logging.getLogger('kubeops')


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    operation = models.CharField(max_length=128, blank=False, null=False)
    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)
    params = common_models.JsonDictTextField(default={})
    steps = common_models.JsonListTextField(default=[], null=True)

    STEP_STATUS_PENDING = 'pending'
    STEP_STATUS_RUNNING = 'running'
    STEP_STATUS_SUCCESS = 'success'
    STEP_STATUS_ERROR = 'error'

    @property
    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        cluster = self.get_cluster()
        settings = Setting.get_db_settings()
        extra_vars = {
            "cluster_name": cluster.name,
            "cluster_domain": cluster.cluster_doamin_suffix
        }
        extra_vars.update(settings)
        extra_vars.update(cluster.configs)
        ignore_errors = False
        return_running = False
        message_client = MessageClient()
        message = {
            "item_id": cluster.item_id,
            "title": self.get_operation_name(),
            "content": "",
            "level": "INFO",
            "type": "SYSTEM"
        }
        try:
            if self.operation == "install":
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                cluster.change_status(Cluster.CLUSTER_STATUS_INSTALLING)
                result = self.on_install(extra_vars)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'uninstall':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                cluster.change_status(Cluster.CLUSTER_STATUS_DELETING)
                result = self.on_uninstall(extra_vars)
                cluster.change_status(Cluster.CLUSTER_STATUS_READY)
                kubeops_api.cluster_monitor.delete_cluster_redis_data(cluster.name)
            elif self.operation == 'bigip-config':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                ignore_errors = True
                result = self.on_f5_config(extra_vars)
            elif self.operation == 'upgrade':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                cluster.change_status(Cluster.CLUSTER_STATUS_UPGRADING)
                package_name = self.params.get('package', None)
                package = Package.objects.get(name=package_name)
                extra_vars.update(package.meta.get('vars'))
                result = self.on_upgrade(extra_vars)
                if result.get('summary', {}).get('success', False):
                    cluster.upgrade_package(package_name)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'scale':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                ignore_errors = True
                return_running = True
                cluster.change_status(Cluster.CLUSTER_DEPLOY_TYPE_SCALING)
                result = self.on_scaling(extra_vars)
                cluster.exit_new_node()
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'add-worker':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                ignore_errors = True
                return_running = True
                cluster.change_status(Cluster.CLUSTER_DEPLOY_TYPE_SCALING)
                result = self.on_add_worker(extra_vars)
                cluster.exit_new_node()
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'remove-worker':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                ignore_errors = True
                return_running = True
                cluster.change_status(Cluster.CLUSTER_DEPLOY_TYPE_SCALING)
                result = self.on_remove_worker(extra_vars)
                if not result.get('summary', {}).get('success', False):
                    cluster.exit_new_node()
                else:
                    node_names = self.params.get('nodes', None)
                    cluster.change_to()
                    nodes = Node.objects.filter(name__in=node_names)
                    for node in nodes:
                        node.delete()
                    cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'restore':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                cluster.change_status(Cluster.CLUSTER_STATUS_RESTORING)
                cluster_backup_id = self.params.get('clusterBackupId', None)
                result = self.on_restore(extra_vars, cluster_backup_id)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            elif self.operation == 'backup':
                logger.info(msg="cluster: {} exec: {} ".format(cluster, self.operation))
                cluster.change_status(Cluster.CLUSTER_STATUS_BACKUP)
                cluster_storage_id = self.params.get('backupStorageId', None)
                result = self.on_backup(extra_vars)
                self.on_upload_backup_file(cluster_storage_id)
                cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
            if not result.get('summary', {}).get('success', False):
                message['content'] = self.get_content(False)
                message['level'] = 'WARNING'
                if not ignore_errors:
                    cluster.change_status(Cluster.CLUSTER_STATUS_ERROR)
                if return_running:
                    cluster.change_status(Cluster.CLUSTER_STATUS_RUNNING)
                logger.error(msg=":cluster {} exec {} error".format(cluster, self.operation), exc_info=True)
            else:
                message['content'] = self.get_content(True)
            message_client.insert_message(message)
        except Exception as e:
            logger.error(msg=":cluster {} exec {} error".format(cluster, self.operation), exc_info=True)
            cluster.change_status(Cluster.CLUSTER_STATUS_ERROR)
            message['content'] = self.get_content(False)
            message['level'] = 'WARNING'
            message_client.insert_message(message)
        post_deploy_execution_start.send(self.__class__, execution=self, result=result, ignore_errors=ignore_errors)
        return result

    def on_install(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('install')
        self.set_step_default()
        self.update_current_step('create-resource', DeployExecution.STEP_STATUS_RUNNING)
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            try:
                cluster.create_resource()
                cluster.refresh_from_db()
                extra_vars.update(cluster.configs)
                self.update_current_step('create-resource', DeployExecution.STEP_STATUS_SUCCESS)
            except RuntimeError as e:
                self.update_current_step('create-resource', DeployExecution.STEP_STATUS_ERROR)
                raise e
        else:
            delete = None
            for step in self.steps:
                if step['name'] == 'create-resource':
                    delete = step
            self.steps.remove(delete)
        return self.run_playbooks(extra_vars)

    def on_scaling(self, extra_vars):
        cluster = self.get_cluster()
        cluster.change_to()
        if not Role.objects.filter(name='new_node'):
            Role.objects.create(name='new_node', project=cluster)
        self.steps = cluster.get_steps('scale')
        self.set_step_default()
        self.update_current_step('create-resource', DeployExecution.STEP_STATUS_RUNNING)
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            try:
                num = self.params.get('num', None)
                cluster.scale_up_to(int(num))
                self.update_current_step('create-resource', DeployExecution.STEP_STATUS_SUCCESS)
            except RuntimeError as e:
                self.update_current_step('create-resource', DeployExecution.STEP_STATUS_ERROR)
                raise e
        return self.run_playbooks(extra_vars)

    def on_add_worker(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('add-worker')
        self.set_step_default()
        host_names = self.params.get('hosts', None)
        hosts = Host.objects.filter(name__in=host_names)
        cluster.add_worker(hosts)
        return self.run_playbooks(extra_vars)

    def on_remove_worker(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('remove-worker')
        self.set_step_default()
        node_names = self.params.get('nodes', None)
        cluster.change_to()
        nodes = Node.objects.filter(name__in=node_names)
        for node in nodes:
            node.set_groups(['new_node', 'worker'])
        return self.run_playbooks(extra_vars)

    def on_uninstall(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('uninstall')
        self.set_step_default()
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            try:
                self.update_current_step('uninstall', DeployExecution.STEP_STATUS_RUNNING)
                cluster.destroy_resource()
                self.update_current_step('uninstall', DeployExecution.STEP_STATUS_SUCCESS)
            except RuntimeError as e:
                self.update_current_step('uninstall', DeployExecution.STEP_STATUS_ERROR)
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

    def on_backup(self, extra_vars):
        cluster = self.get_cluster()
        self.steps = cluster.get_steps('cluster-backup')
        return self.run_playbooks(extra_vars)

    def on_upload_backup_file(self, backup_storage_id):
        cluster = self.get_cluster()
        return kubeops_api.cluster_backup_utils.upload_backup_file(cluster.id, backup_storage_id)

    def run_playbooks(self, extra_vars):
        result = {"raw": {}, "summary": {}}
        for step in self.steps:
            playbook_name = step.get('playbook', None)
            if playbook_name:
                playbook = self.project.playbook_set.get(name=playbook_name)
                self.update_current_step(step['name'], DeployExecution.STEP_STATUS_RUNNING)
                _result = playbook.execute(extra_vars=extra_vars)
                result["summary"].update(_result["summary"])
                self.update_current_step(step['name'], DeployExecution.STEP_STATUS_SUCCESS)
                if not _result.get('summary', {}).get('success', False):
                    self.update_current_step(step['name'], DeployExecution.STEP_STATUS_ERROR)
                    return result
        return result

    def set_step_default(self):
        for step in self.steps:
            step['status'] = DeployExecution.STEP_STATUS_PENDING

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

    def mark_state(self, state):
        self.state = state
        self.date_end = timezone.now()
        self.timedelta = (timezone.now() - self.date_start).seconds
        self.save()

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)

    def get_operation_name(self):
        operation_name = {
            "install": "集群安装",
            "uninstall": "集群卸载",
            "upgrade": "集群升级",
            "scale": "集群伸缩",
            "add-worker": "集群伸缩",
            "remove-worker": "集群安装",
            "restore": "集群恢复",
            "backup": "集群备份",
        }
        return operation_name[self.operation]

    def get_content(self, success):
        cluster = self.get_cluster()
        content = {
            "item_name": cluster.item_name,
            "resource": "集群",
            "resource_name": cluster.name,
            "resource_type": 'CLUSTER',
            "detail": self.get_msg_detail(success),
            "status": cluster.status
        }
        return content

    def get_msg_detail(self, success):
        operation = self.get_operation_name()
        if success:
            result = "成功"
        else:
            result = "失败"
        return json.dumps({"message": operation + result})
