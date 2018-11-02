# -*- coding: utf-8 -*-
#
from django.dispatch import Signal


pre_playbook_exec = Signal(providing_args=('playbook', 'save_history'))
post_playbook_exec = Signal(providing_args=('playbook', 'save_history', 'result'))

pre_adhoc_exec = Signal(providing_args=('adhoc', 'save_history'))
post_adhoc_exec = Signal(providing_args=('adhoc', 'save_history', 'result'))
