import os

from django.db.models.signals import m2m_changed, post_save, pre_save
from django.dispatch import receiver
from django.utils import timezone
from gunicorn.config import Setting

from openshift_api.models.cluster import Cluster
from openshift_api.models.host import HostInfo
from openshift_api.models.node import Node
from openshift_api.models.package import Package
from openshift_api.models.storage import Storage
from .signals import pre_deploy_execution_start, post_deploy_execution_start


@receiver(post_save, sender=Cluster)
def on_cluster_save(sender, instance=None, created=True, **kwargs):
    if created and instance and instance.template:
        instance.on_cluster_create()


@receiver(post_save, sender=Storage)
def on_storage_save(sender, instance=None, **kwargs):
    if instance and instance.template:
        instance.on_storage_create()


@receiver(post_save, sender=Node)
def on_node_save(sender, instance=None, created=False, **kwargs):
    if created and not instance.name == 'localhost' and not instance.name == '127.0.0.1':
        instance.on_node_save()


@receiver(pre_save, sender=HostInfo)
def before_hostInfo_save(sender, instance=None, created=False, **kwargs):
    if created:
        instance.clear_host_info()


@receiver(pre_save, sender=Setting)
def before_setting_save(sender, instance=None, **kwargs):
    pass
    if instance.name == 'hostname':
        os.putenv("REGISTORY_HOSTNAME", instance.value)


def auto_lookup_packages():
    try:
        Package.lookup()
    except:
        pass


@receiver(pre_deploy_execution_start)
def on_execution_start(sender, execution, **kwargs):
    execution.date_start = timezone.now()
    execution.state = execution.STATE_STARTED
    execution.save()


@receiver(post_deploy_execution_start)
def on_execution_end(sender, execution, result, **kwargs):
    date_finished = timezone.now()
    timedelta = (timezone.now() - execution.date_start).seconds
    state = execution.STATE_FAILURE
    if result.get('summary', {}).get("success", False):
        state = execution.STATE_SUCCESS
    execution.result_summary = result.get('summary', {})
    execution.result_raw = result.get('raw', {})
    execution.state = state
    execution.date_end = date_finished
    execution.timedelta = timedelta
    execution.save()


auto_lookup_packages()
