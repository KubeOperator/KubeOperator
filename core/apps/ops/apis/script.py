from django.db import transaction
from rest_framework.viewsets import ModelViewSet

from kubeops_api.mixin import PageModelViewSet
from ops.models import Script
from ops.models.script import ScriptExecution
from ops.serializers.script import ScriptSerializer
from ops.tasks import start_script_execution


class ScriptViewSet(PageModelViewSet):
    queryset = Script.objects.all()
    serializer_class = ScriptSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class ScriptExecutionViewSet(ModelViewSet):
    queryset = ScriptExecution.objects.all()

    def perform_create(self, serializer):
        instance = serializer.save()
        transaction.on_commit(lambda: start_script_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        ))
        return instance
