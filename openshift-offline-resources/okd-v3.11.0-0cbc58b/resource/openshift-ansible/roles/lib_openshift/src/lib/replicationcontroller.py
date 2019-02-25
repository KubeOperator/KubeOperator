# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-public-methods
class ReplicationController(DeploymentConfig):
    ''' Class to model a replicationcontroller openshift object.

        Currently we are modeled after a deployment config since they
        are very similar.  In the future, when the need arises we
        will add functionality to this class.
    '''
    replicas_path = "spec.replicas"
    env_path = "spec.template.spec.containers[0].env"
    volumes_path = "spec.template.spec.volumes"
    container_path = "spec.template.spec.containers"
    volume_mounts_path = "spec.template.spec.containers[0].volumeMounts"

    def __init__(self, content):
        ''' Constructor for ReplicationController '''
        super(ReplicationController, self).__init__(content=content)
