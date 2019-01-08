# -*- coding: utf-8 -*-
#
from django.dispatch import Signal

pre_deploy_execution_start = Signal(providing_args=('execution',))
post_deploy_execution_start = Signal(providing_args=('execution', 'result'))
