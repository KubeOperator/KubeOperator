from django.db import models
from rest_framework.generics import get_object_or_404
from rest_framework.mixins import ListModelMixin

from ansible_api.api.mixin import ProjectResourceAPIMixin
from kubeops_api.models import Item, ItemResource


class ClusterResourceAPIMixin(ProjectResourceAPIMixin):
    lookup_kwargs = 'cluster_name'


