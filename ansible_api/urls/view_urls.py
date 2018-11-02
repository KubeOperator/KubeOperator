# ~*~ coding: utf-8 ~*~
from django.urls import path

from .. import views


__all__ = ["urlpatterns"]

app_name = "ansible_api"

urlpatterns = [
    # path('celery/<uuid:task_id>/log/', views.celery_log_view),
]
