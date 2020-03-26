import logging

from celery import shared_task

from users.sync.ldap import LDAPSync

logger = logging.getLogger("kubeops")


@shared_task
def start_sync_user_form_ldap():
    sync = LDAPSync()
    sync.run()
