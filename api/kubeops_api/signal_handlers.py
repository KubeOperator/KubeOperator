import json
import os

from django.db.models.signals import post_save, pre_save, post_delete
from django.dispatch import receiver
from django.utils import timezone

from kubeops_api.adhoc import test_host
from kubeops_api.models import Credential
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.package import Package
from .signals import pre_deploy_execution_start, post_deploy_execution_start


@receiver(post_save, sender=Cluster)
def on_cluster_save(sender, instance=None, created=True, **kwargs):
    if created and instance and instance.template:
        instance.on_cluster_create()


@receiver(post_delete, sender=Cluster)
def on_cluster_delete(sender, instance=None, **kwargs):
    instance.on_cluster_delete()


@receiver(post_save, sender=Node)
def on_node_save(sender, instance=None, created=False, **kwargs):
    if created and not instance.name == 'localhost' and not instance.name == '127.0.0.1' and not instance.name == '::1':
        instance.on_node_save()


@receiver(post_save, sender=Host)
def post_host_save(sender, instance=None, created=False, **kwargs):
    if created and instance.auto_gather_info:
        instance.full_host_credential()
        instance.gather_info()


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
def on_execution_end(sender, execution, result, ignore_errors, **kwargs):
    cluster = Cluster.objects.get(id=execution.project.id)
    date_finished = timezone.now()
    timedelta = (timezone.now() - execution.date_start).seconds
    if result.get('summary', {}).get("success", False):
        state = execution.STATE_SUCCESS
    else:
        state = execution.STATE_FAILURE
        if not ignore_errors:
            cluster.change_status(Cluster.CLUSTER_STATUS_ERROR)
    execution.result_summary = result.get('summary', {})
    execution.result_raw = result.get('raw', {})
    execution.state = state
    execution.date_end = date_finished
    execution.timedelta = timedelta
    execution.save()


auto_lookup_packages()
