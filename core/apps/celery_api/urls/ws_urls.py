from django.urls import path

from .. import ws

app_name = 'celery-api'

urlpatterns = [
    path('ws/tasks/<uuid:task_id>/log/', ws.CeleryLogWebsocket, name='task-log-ws'),
]
