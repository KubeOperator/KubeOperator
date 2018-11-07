from django.db.models.signals import m2m_changed
from django.dispatch import receiver

from ansible_api.models import Group


@receiver(m2m_changed, sender=Group.hosts.through)
def on_role_hosts_change(sender, action, instance, reverse, model, pk_set, **kwargs):
    nodes = model.objects.filter(pk__in=pk_set)
    if action == "post_remove":
        instance.on_nodes_leave(nodes)
    elif action == "post_add":
        instance.on_nodes_join(nodes)
    else:
        print("What ever: {}".format(action))
