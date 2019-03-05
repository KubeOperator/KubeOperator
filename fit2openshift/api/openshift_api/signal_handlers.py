import os

from django.db.models.signals import m2m_changed, post_save, pre_save
from django.dispatch import receiver
from django.utils import timezone

from .signals import pre_deploy_execution_start, post_deploy_execution_start
from .models import Role, Package, Cluster, Node, Host, Setting


@receiver(post_save, sender=Cluster)
def on_cluster_save(sender, instance=None, **kwargs):
    if instance and instance.template:
        instance.on_cluster_create()


@receiver(post_save, sender=Host)
def on_host_save(sender, instance=None, created=False, **kwargs):
    if created:
        instance.get_host_info()


@receiver(post_save, sender=Node)
def on_node_save(sender, instance=None, created=False, **kwargs):
    if created and not instance.name == 'localhost':
        instance.on_node_save()


@receiver(pre_save, sender=Node)
def before_node_save(sender, instance=None, created=False, **kwargs):
    if created and not instance.name == 'localhost':
        instance.before_node_save()


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
    ##检查是否存在有效的hostname
    hostname = Setting.objects.filter(key="hostname").first()
    if hostname and hostname.value and not hostname.value == "":
        execution.date_start = timezone.now()
        execution.state = execution.STATE_STARTED
        execution.save()
    else:
        raise Exception("Hostname not found!")


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
