from django.urls import path

from .. import api

app_name = 'celery-api'

urlpatterns = [
    path('tasks/<uuid:pk>/result/', api.TaskResultApi.as_view(), name='task-result-api'),
    path('tasks/<uuid:pk>/log/', api.TaskLogApi.as_view(), name='task-log-api'),
]
