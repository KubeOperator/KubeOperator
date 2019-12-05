from rest_framework import viewsets

from ansible_api.permissions import IsSuperUser
from storage import serializers
from storage.models import NfsStorage,CephStorage


class NfsStorageViewSet(viewsets.ModelViewSet):
    queryset = NfsStorage.objects.all()
    serializer_class = serializers.NfsStorageSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

class CephStorageViewSet(viewsets.ModelViewSet):
    queryset = CephStorage.objects.all()
    serializer_class = serializers.CephStorageSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'