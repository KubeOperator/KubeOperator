import os
from http import HTTPStatus

from django.http import JsonResponse
from rest_framework.generics import get_object_or_404, CreateAPIView
from kubeoperator.settings import MEDIA_DIR
from kubeops_api.host_import import HostImporter
from kubeops_api.mixin import PageModelViewSet
from kubeops_api.models.host import Host
from kubeops_api.serializers.host import HostSerializer
from kubeops_api.models.item import Item
from kubeops_api.models.item_resource import ItemResource

__all__ = ["HostViewSet"]


class HostImportAPIView(CreateAPIView):

    def create(self, request, *args, **kwargs):
        source = request.data.get("source", None)
        for item in source:
            importer = HostImporter(path=os.path.join(MEDIA_DIR, item))
            importer.run()
        return JsonResponse(data={"success": True}, status=HTTPStatus.CREATED)


class HostViewSet(PageModelViewSet):
    queryset = Host.objects.all()
    serializer_class = HostSerializer

    def retrieve(self, request, *args, **kwargs):
        pk = kwargs.get('pk')
        host = get_object_or_404(Host, pk=pk)
        host.gather_info(retry=1)
        return super().retrieve(request, *args, **kwargs)

    def list(self, request, *args, **kwargs):
        item_name = request.query_params.get('item', None)
        if item_name:
            item = get_object_or_404(Item, name=item_name)
            resources = ItemResource.objects.filter(item_id=item.id)
            if resources:
                self.queryset = Host.objects.filter(id__in=resources.values_list("resource_id"))
        return super().list(self, request, *args, **kwargs)
