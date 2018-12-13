# -*- coding: utf-8 -*-
#
from django.dispatch import Signal


pre_execution_start = Signal(providing_args=('execution',))
post_execution_start = Signal(providing_args=('execution', 'result'))

pre_adhoc_start = Signal(providing_args=('execution', 'save_history'))
post_adhoc_start = Signal(providing_args=('execution', 'save_history', 'result'))
