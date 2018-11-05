# -*- coding: utf-8 -*-
#

from rest_framework import viewsets
from rest_framework.response import Response
from django.shortcuts import get_object_or_404

from common.api import LogTailApi
from .mixin import ProjectResourceAPIMixin
from ..permissions import IsSuperUser
from ..serializers import (
    PlaybookReadSerializer, PlaybookSerializer,
    PlaybookExecutionSerializer, TaskSerializer,
)
from ..models import Playbook, PlaybookExecution
from ..tasks import execute_playbook, start_playbook_execution


__all__ = [
    'ProjectPlaybookViewSet',
    'PlaybookExecutionViewSet',
]


class ProjectPlaybookViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Playbook.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = PlaybookSerializer
    read_serializer_class = PlaybookReadSerializer


class PlaybookExecutionViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = PlaybookExecution.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = PlaybookExecutionSerializer
    http_method_names = ['post', 'get', 'option', 'head']

    def perform_create(self, serializer):
        instance = serializer.save()
        start_playbook_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        )
        return instance
