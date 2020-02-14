# ~*~ coding: utf-8 ~*~
from rest_framework import viewsets

from .mixin import ProjectResourceAPIMixin
from ..permissions import IsSuperUser
from ..models import AdHoc, AdHocExecution
from ..serializers import AdHocReadSerializer, AdHocSerializer, AdHocExecutionSerializer
from ..tasks import start_adhoc_execution


__all__ = [
    'AdHocViewSet', 'AdHocExecutionViewSet',
]


class AdHocViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = AdHoc.objects.all()

    serializer_class = AdHocSerializer
    read_serializer_class = AdHocReadSerializer


class AdHocExecutionViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = AdHocExecution.objects.all()

    serializer_class = AdHocExecutionSerializer

    http_method_names = ['post', 'get', 'option', 'head']

    def perform_create(self, serializer):
        instance = serializer.save()
        start_adhoc_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        )
        return instance
