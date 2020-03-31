from rest_framework import viewsets
from ansible_api.api.mixin import ProjectResourceAPIMixin
from common.api import Pagination


class ClusterResourceAPIMixin(ProjectResourceAPIMixin):
    lookup_kwargs = 'cluster_name'


class PageModelViewSet(viewsets.ModelViewSet):
    pagination_class = Pagination

    def list(self, request, *args, **kwargs):
        if not self.request.query_params.get('page', None):
            self.pagination_class = None
        return super().list(request, *args, **kwargs)
