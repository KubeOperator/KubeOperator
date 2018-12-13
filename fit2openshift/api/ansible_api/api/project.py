# -*- coding:utf-8 -*-
#

from rest_framework import viewsets
from django.db import transaction

from ..permissions import IsSuperUser
from ..serializers import ProjectSerializer
from ..models import Project

__all__ = [
    'ProjectViewSet',
]


class ProjectViewSet(viewsets.ModelViewSet):
    queryset = Project.objects.all()
    permission_classes = (IsSuperUser,)
    filter_fields = ('name',)
    serializer_class = ProjectSerializer
    lookup_url_kwarg = 'name'
    lookup_field = 'name'

    @transaction.atomic
    def dispatch(self, request, *args, **kwargs):
        return super().dispatch(request, *args, **kwargs)





