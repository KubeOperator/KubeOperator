# -*- coding: utf-8 -*-
#

from django.http import Http404
from rest_framework import viewsets
from rest_framework.generics import RetrieveUpdateAPIView

from .mixin import ProjectResourceAPIMixin
from ..permissions import IsSuperUser
from ..models import ClusterHost, ClusterGroup, Group, Host
from ..serializers import (
    ClusterHostSerializer, ClusterGroupSerializer, HostSerializer, GroupSerializer,
    InventorySerializer,
)

__all__ = [
    'HostViewSet', 'GroupViewSet', 'ProjectHostViewSet',
    'ProjectGroupViewSet',
    'ProjectInventoryApi',
]


class HostViewSet(viewsets.ModelViewSet):
    serializer_class = ClusterHostSerializer

    queryset = ClusterHost.objects.all()


class GroupViewSet(viewsets.ModelViewSet):
    serializer_class = ClusterGroupSerializer
    queryset = ClusterGroup.objects.all()



class ProjectHostViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    serializer_class = HostSerializer

    queryset = Host.objects.all()


class ProjectGroupViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    serializer_class = GroupSerializer
    queryset = Group.objects.all()



class ProjectInventoryApi(ProjectResourceAPIMixin, RetrieveUpdateAPIView):
    serializer_class = InventorySerializer


    def get_object(self):
        return self.project.inventory
