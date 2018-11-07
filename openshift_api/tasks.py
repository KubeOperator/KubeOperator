from celery import shared_task


@shared_task
def start_openshift_deploy(cluster_id):
    from .models import Cluster
    cluster = Cluster.objects.get(id=cluster_id)
    cluster.change_to()
    return cluster.execute()
