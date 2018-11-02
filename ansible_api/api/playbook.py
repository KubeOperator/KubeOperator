# -*- coding: utf-8 -*-
#

from rest_framework import viewsets
from rest_framework.response import Response
from django.shortcuts import get_object_or_404

from .mixin import ProjectResourceAPIMixin
from ..permissions import IsSuperUser
from ..serializers import (
    PlaybookReadSerializer, PlaybookSerializer,
    PlaybookExecutionSerializer, TaskSerializer,
)
from ..models import Playbook, PlaybookExecution
from ..tasks import execute_playbook


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


