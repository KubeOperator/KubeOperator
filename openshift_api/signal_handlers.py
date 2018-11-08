from django.db.models.signals import m2m_changed
from django.dispatch import receiver

from ansible_api.models import Group
from .models import Role


@receiver(m2m_changed, sender=Group.hosts.through)
def on_role_hosts_change(sender, action, instance, reverse, model, pk_set, **kwargs):
    Role.update_node_group_labels()
