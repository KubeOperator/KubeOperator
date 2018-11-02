# ~*~ coding: utf-8 ~*~
from rest_framework import viewsets, generics
from rest_framework.pagination import LimitOffsetPagination
from rest_framework.response import Response

from common.api import LogTailApi
from .mixin import ProjectResourceAPIMixin
from ..permissions import IsSuperUser
from ..models import AdHoc, AdHocExecution
from ..serializers import AdHocReadSerializer, AdHocSerializer, AdHocExecutionSerializer
from ..tasks import execute_adhoc


__all__ = [
    'AdHocViewSet', 'AdHocExecutionViewSet', 'AdHocLogApi',
]


class AdHocViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = AdHoc.objects.all()
    permission_classes = (IsSuperUser,)
    pagination_class = LimitOffsetPagination
    serializer_class = AdHocSerializer


class AdHocExecutionViewSet(ProjectResourceAPIMixin, viewsets.ModelViewSet):
    queryset = AdHocExecution.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = AdHocExecutionSerializer

    http_method_names = ['post', 'get', 'option', 'head']

    # def update(self, request, *args, **kwargs):
    #     # serializer = self.get_serializer(data=request.data)
    #     # serializer.is_valid(raise_exception=True)
    #     adhoc = self.get_object()
    #     # passwords = serializer.validated_data.get('passwords')
    #     task = execute_adhoc.delay(str(adhoc.id))
    #     return Response({"task": task.id})
    #
    # def retrieve(self, request, *args, **kwargs):
    #     return self.update(request, *args, **kwargs)


class AdHocLogApi(LogTailApi):
    queryset = AdHoc.objects.all()
    object = None

    def get(self, request, *args, **kwargs):
        self.object = self.get_object()
        return super().get(request, *args, **kwargs)

    def get_log_path(self):
        return self.object.log_path

    def is_end(self):
        return self.object.is_finished
