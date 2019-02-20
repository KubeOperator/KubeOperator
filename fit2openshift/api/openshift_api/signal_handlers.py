from django.db.models.signals import m2m_changed, post_save, pre_save
from django.dispatch import receiver
from django.utils import timezone

from .signals import pre_deploy_execution_start, post_deploy_execution_start
from .models import Role, Package, Cluster, Node, Host


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
    execution.date_finished = date_finished
    execution.timedelta = timedelta
    execution.save()


auto_lookup_packages()
