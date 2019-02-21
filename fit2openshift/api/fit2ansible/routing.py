from channels.auth import AuthMiddlewareStack
from channels.routing import ProtocolTypeRouter, URLRouter
from django.conf.urls import url
from django.urls import path

from celery_api.ws import CeleryLogWebsocket
from openshift_api.ws import F2OWebsocket

application = ProtocolTypeRouter({
    # Empty for now (http->django views is added by default)
    'websocket': AuthMiddlewareStack(
        URLRouter([
            path('ws/tasks/<uuid:task_id>/log/', CeleryLogWebsocket, name='task-log-ws'),
            path('ws/progress/<uuid:execution_id>/', F2OWebsocket, name='execution-progress-ws'),
        ]),
    ),

})
