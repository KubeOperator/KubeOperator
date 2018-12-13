from ansible_api.api.mixin import ProjectResourceAPIMixin


class ClusterResourceAPIMixin(ProjectResourceAPIMixin):
    lookup_kwargs = 'cluster_name'
