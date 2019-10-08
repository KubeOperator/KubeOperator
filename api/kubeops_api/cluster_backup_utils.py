
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
                run_backup(project_id,backup_storage_id)
        else:
            run_backup(project_id,backup_storage_id)

def run_backup(project_id,backup_storage_id):
    playbook_name = 'cluster-backup'
    cluster = Cluster.objects.get(id=project_id)
    hostname = Setting.objects.get(key='local_hostname')
    domain_suffix = Setting.objects.get(key="domain_suffix")
    project = Project.objects.get(id=project_id)
    playbook = project.playbook_set.get(name=playbook_name)
    backup_storage = BackupStorage.objects.get(id=backup_storage_id)
    extra_vars = {
        "cluster_name": cluster.name,
        "local_hostname": hostname.value,
        "domain_suffix": domain_suffix.value,
        "APP_DOMAIN": "apps.{}.{}".format(cluster.name, domain_suffix.value),
        "NODE_NAME": cluster.name
    }
    extra_vars.update(cluster.configs)
    _result = playbook.execute(extra_vars=extra_vars)
    now = datetime.datetime.now().strftime('%Y-%m-%d')
    client = StorageClient(backup_storage)
    file_name = cluster.name+str(now)+'.zip'
    file_remote_path = cluster.name+'/'+file_name
    result,message = client.upload_file("/etc/ansible/roles/cluster-backup/files/cluster-backup.zip",file_remote_path)
    if result:
        clusterBackup = ClusterBackup(name=file_name,size=10,folder=file_remote_path,
                                      backup_storage_id=backup_storage_id,project_id=project_id)
        clusterBackup.save()
    else:
        pass


def run_restore(cluster_backup_id):
    playbook_name = 'cluster-restore'
    cluster_backup = ClusterBackup.objects.get(id=cluster_backup_id)
    backup_storage = BackupStorage.objects.get(id=cluster_backup.backup_storage_id)
    project = Project.objects.get(id=cluster_backup.project_id)
    playbook = project.playbook_set.get(name=playbook_name)
    cluster = Cluster.objects.get(id=cluster_backup.project_id)
    client = StorageClient(backup_storage)
    backup_file_path = cluster.name+'/'+cluster_backup.name
    if client.exists(backup_file_path):
        success = client.download_file(backup_file_path,"/etc/ansible/roles/cluster-backup/files/cluster-backup.zip")
        if success:
            extra_vars = {
                "cluster_name": cluster.name,
                "NODE_NAME": cluster.name
            }
            _result = playbook.execute(extra_vars=extra_vars)
            if _result:
                return True
            else:
                return False
        else:
            raise Exception('download file failed!')

    else:
        raise Exception('File is not exist!')


