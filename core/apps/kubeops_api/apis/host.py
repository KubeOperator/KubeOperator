from rest_framework.generics import get_object_or_404
from kubeops_api.models.host import Host
from kubeops_api.serializers.host import HostSerializer
from rest_framework import viewsets
from kubeops_api.models.item import Item
from kubeops_api.models.item_resource import ItemResource

__all__ = ["HostViewSet"]


class HostViewSet(viewsets.ModelViewSet):
    queryset = Host.objects.all()
    serializer_class = HostSerializer

    def retrieve(self, request, *args, **kwargs):
        pk = kwargs.get('pk')
        host = get_object_or_404(Host, pk=pk)
        host.gather_info(retry=1)
        return super().retrieve(request, *args, **kwargs)

    def list(self, request, *args, **kwargs):
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            item = Item.objects.get(name=itemName)
            resource_ids = ItemResource.objects.filter(item_id=item.id).values_list("resource_id")
            self.queryset = Host.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        else:
            return super().list(self, request, *args, **kwargs)