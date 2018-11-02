# -*- coding: utf-8 -*-
#

from rest_framework import viewsets

from .mixin import ProjectObjectMixin
from ..permissions import IsSuperUser
from ..models import Role
from ..serializers import RoleReadSerializer, RoleSerializer


__all__ = [
    'RoleViewSet', 'ProjectRoleViewSet',
]


class RoleViewSet(viewsets.ModelViewSet):
    queryset = Role.objects.all()
    serializer_class = RoleReadSerializer
    permission_classes = (IsSuperUser,)


class ProjectRoleViewSet(ProjectObjectMixin, RoleViewSet):
    serializer_class = RoleSerializer
