from rest_framework.generics import get_object_or_404
from kubeops_api.models.host import Host
from kubeops_api.serializers.host import HostSerializer
from rest_framework import viewsets

__all__ = ["HostViewSet"]


class HostViewSet(viewsets.ModelViewSet):
    queryset = Host.objects.all()
    serializer_class = HostSerializer

    def retrieve(self, request, *args, **kwargs):
        pk = kwargs.get('pk')
        host = get_object_or_404(Host, pk=pk)
        host.gather_info(retry=1)
        return super().retrieve(request, *args, **kwargs)
