from celery import shared_task

from common.utils import get_object_or_none
<<<<<<< HEAD
from ansible_api.ctx import change_to_root
=======
from ansible_api.ctx import change_to_root, set_current_project
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce

from .models import DeployExecution


@shared_task
<<<<<<< HEAD
=======
def start_openshift_deploy(cluster_id):
    from .models import Cluster
    cluster = Cluster.objects.get(id=cluster_id)
    cluster.change_to()
    return cluster.execute()


@shared_task
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
def start_deploy_execution(eid, **kwargs):
    change_to_root()
    execution = get_object_or_none(DeployExecution, id=eid)
    if execution:
<<<<<<< HEAD
        execution.project.change_to()
=======
        set_current_project(execution.project)
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
        return execution.start()
    else:
        msg = "No execution found: {}".format(eid)
        print(msg)
        return {"error": msg}
