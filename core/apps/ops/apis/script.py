from kubeops_api.mixin import PageModelViewSet
from ops.models import Script
from ops.serializers.script import ScriptSerializer


class ScriptViewSet(PageModelViewSet):
    queryset = Script.objects.all()
    serializer_class = ScriptSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'
