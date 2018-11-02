from django.urls import path

from .. import ws

urlpatterns = [
    path('ws/celery/<uuid:task_id>/log/', ws.CeleryLogWebsocket),
]
