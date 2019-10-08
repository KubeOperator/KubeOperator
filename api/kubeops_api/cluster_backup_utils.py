
from kubeops_api.models.cluster_backup import ClusterBackup
from kubeops_api.models.backup_strategy import BackupStrategy
from ansible_api.models.project import Project
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.setting import Setting
from kubeops_api.models.backup_storage import BackupStorage
from kubeops_api.storage_client import StorageClient
import datetime
import time


def cluster_backup():
    backup_strategies = BackupStrategy.objects.all()
    for b in backup_strategies:
        project_id = str(b.project_id)
        backup_storage_id = str(b.backup_storage_id)
        cluster_backups = ClusterBackup.objects.filter(backup_storage_id=backup_storage_id).order_by('date_created')
        if cluster_backups:
            time_now = datetime.datetime.now().strftime('%Y-%m-%d')
            date1 = time.strptime(time_now,'%Y-%m-%d')
            time_backup = cluster_backups[0].date_created.strftime('%Y-%m-%d')
            date2 = time.strptime(time_backup,'%Y-%m-%d')
            d1 = datetime.datetime(date1[0], date1[1], date1[2])
            d2 = datetime.datetime(date2[0], date2[1], date2[2])
            day = (d1-d2).days
            if day >= b.cron:
                success = run_backup(project_id,backup_storage_id)
                if success and len (cluster_backups)+1 > b.save_num:
                    delete_backup(cluster_backups[-1].id)
        else:
            run_backup(project_id,backup_storage_id)

def run_backup(project_id,backup_storage_id):
    cluster = Cluster.objects.get(id=project_id)
    steps = cluster.get_steps('cluster-backup')
    cluster.configs = cluster.load_config_file()
    hostname = Setting.objects.get(key='local_hostname')
    domain_suffix = Setting.objects.get(key="domain_suffix")
    project = Project.objects.get(id=project_id)
    print(project)
    backup_storage = BackupStorage.objects.get(id=backup_storage_id)
    extra_vars = {
        "cluster_name": cluster.name,
        "local_hostname": hostname.value,
        "domain_suffix": domain_suffix.value,
        "APP_DOMAIN": "apps.{}.{}".format(cluster.name, domain_suffix.value),
    }
    extra_vars.update(cluster.configs)
    run_playbooks(steps,extra_vars,project)
    now = datetime.datetime.now().strftime('%Y-%m-%d')
    client = StorageClient(backup_storage)
    client.check_valid()
    file_name = cluster.name+'-'+str(now)+'.zip'
    file_remote_path = cluster.name+'/'+file_name
    result,message = client.upload_file("/etc/ansible/roles/cluster-backup/files/cluster-backup.zip",file_remote_path)
    if result:
        clusterBackup = ClusterBackup(name=file_name,size=10,folder=file_remote_path,
                                      backup_storage_id=backup_storage_id,project_id=project_id)
        clusterBackup.save()
        return True

def run_restore(cluster_backup_id):
    cluster_backup = ClusterBackup.objects.get(id=cluster_backup_id)
    backup_storage = BackupStorage.objects.get(id=cluster_backup.backup_storage_id)
    project = Project.objects.get(id=cluster_backup.project_id)
    cluster = Cluster.objects.get(id=cluster_backup.project_id)
    steps = cluster.get_steps('cluster-backup')
    cluster.configs = cluster.load_config_file()
    client = StorageClient(backup_storage)
    backup_file_path = cluster.name+'/'+cluster_backup.name
    if client.exists(backup_file_path):
        success = client.download_file(backup_file_path,"/etc/ansible/roles/cluster-backup/files/cluster-backup.zip")
        if success:
            extra_vars = {
                "cluster_name": cluster.name,
            }
            run_playbooks(steps, extra_vars, project)
        else:
            raise Exception('download file failed!')
    else:
        raise Exception('File is not exist!')


def delete_backup(cluster_backup_id):
    cluster_backup = ClusterBackup.objects.get(id=cluster_backup_id)
    cluster_backup.delete()

def run_playbooks(steps,extra_vars,project):
    result = {"raw": {}, "summary": {}}
    for step in steps:
        playbook_name = step.get('playbook', None)
        if playbook_name:
            project.change_to()
            playbook = project.playbook_set.get(name=playbook_name)
            _result = playbook.execute(extra_vars=extra_vars)
            result["summary"].update(_result["summary"])
            if not _result.get('summary', {}).get('success', False):
                raise RuntimeError("playbook: {} error!".format(step['playbook']))
    return result
