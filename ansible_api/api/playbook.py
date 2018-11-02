# -*- coding: utf-8 -*-
#

from rest_framework import viewsets
from rest_framework.response import Response

from .mixin import ProjectObjectMixin
from ..permissions import IsSuperUser
from ..serializers import (
    PlaybookReadSerializer, PlaybookSerializer,
    PlaybookExecutionSerializer, TaskSerializer,
)
from ..models import Playbook, PlaybookExecution
from ..tasks import execute_playbook


__all__ = [
    'ProjectPlaybookViewSet',
    'ProjectPlaybookExecutionViewSet',
]


class ProjectPlaybookViewSet(ProjectObjectMixin, viewsets.ModelViewSet):
    queryset = Playbook.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = PlaybookSerializer
    read_serializer_class = PlaybookReadSerializer


class ProjectPlaybookExecutionViewSet(ProjectObjectMixin, viewsets.ModelViewSet):
    queryset = PlaybookExecution.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = PlaybookExecutionSerializer

    def retrieve(self, request, *args, **kwargs):
        return self.update(request, *args, **kwargs)

    def update(self, request, *args, **kwargs):
        instance = self.get_object()
        task = execute_playbook.delay(str(instance.id))
        return Response({"task": task.id})

