# -*- coding: utf-8 -*-
#
from rest_framework.viewsets import ReadOnlyModelViewSet
from rest_framework.response import Response

from ..permissions import AdminUserRequiredMixin

from ..ansible.modules import AnsibleModules
from ..serializers import ModuleSerializer


__all__ = ['AnsibleModuleViewSet']


class AnsibleModuleViewSet(AdminUserRequiredMixin, ReadOnlyModelViewSet):
    serializer_class = ModuleSerializer

    def list(self, request, *args, **kwargs):
        category = AnsibleModules.category_with_modules()
        category.pop('windows', None)
        category.pop('command', None)
        category.pop('inventory_obj', None)
        category.pop('remote_management', None)
        return Response(category)

    def retrieve(self, request, *args, **kwargs):
        name = self.kwargs.get('pk')
        detail = AnsibleModules.modules().get(name)
        return Response(detail)
