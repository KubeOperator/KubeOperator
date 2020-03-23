from django.http import JsonResponse
from rest_framework.generics import RetrieveAPIView, get_object_or_404
from kubeops_api.grade import query_cluster_grade
from kubeops_api.models.cluster import Cluster
from validator.base import ClusterResultJsonEncoder


class GradeRetrieveAPIView(RetrieveAPIView):
    def get(self, request, *args, **kwargs):
        pk = kwargs.get("pk")
        cluster = get_object_or_404(Cluster, pk=pk)
        data = query_cluster_grade(cluster)
        return JsonResponse(status=201, data=data, safe=False, encoder=ClusterResultJsonEncoder)
