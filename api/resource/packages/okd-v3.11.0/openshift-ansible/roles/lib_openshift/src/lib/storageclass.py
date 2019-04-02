# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class StorageClassConfig(object):
    ''' Handle service options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 name,
                 provisioner,
                 parameters=None,
                 annotations=None,
                 default_storage_class="false",
                 api_version='v1',
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 mount_options=None,
                 reclaim_policy=None):
        ''' constructor for handling storageclass options '''
        self.name = name
        self.parameters = parameters
        self.annotations = annotations
        self.provisioner = provisioner
        self.api_version = api_version
        self.default_storage_class = str(default_storage_class).lower()
        self.kubeconfig = kubeconfig
        self.mount_options = mount_options
        self.reclaim_policy = reclaim_policy
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' instantiates a storageclass dict '''
        self.data['apiVersion'] = self.api_version
        self.data['kind'] = 'StorageClass'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name

        self.data['metadata']['annotations'] = {}
        if self.annotations is not None:
            self.data['metadata']['annotations'] = self.annotations

        self.data['metadata']['annotations']['storageclass.beta.kubernetes.io/is-default-class'] = \
                self.default_storage_class

        self.data['provisioner'] = self.provisioner

        self.data['parameters'] = {}
        if self.parameters is not None:
            self.data['parameters'].update(self.parameters)

        # default to aws if no params were passed
        else:
            self.data['parameters']['type'] = 'gp2'

        self.data['mountOptions'] = self.mount_options or []

        if self.reclaim_policy is not None:
            self.data['reclaimPolicy'] = self.reclaim_policy



# pylint: disable=too-many-instance-attributes,too-many-public-methods
class StorageClass(Yedit):
    ''' Class to model the oc storageclass object '''
    annotations_path = "metadata.annotations"
    provisioner_path = "provisioner"
    parameters_path = "parameters"
    mount_options_path = "mountOptions"
    reclaim_policy_path = "reclaimPolicy"
    kind = 'StorageClass'

    def __init__(self, content):
        '''StorageClass constructor'''
        super(StorageClass, self).__init__(content=content)

    def get_annotations(self):
        ''' get a list of ports '''
        return self.get(StorageClass.annotations_path) or {}

    def get_parameters(self):
        ''' get the service selector'''
        return self.get(StorageClass.parameters_path) or {}

    def get_mount_options(self):
        ''' get mount options'''
        return self.get(StorageClass.mount_options_path) or []

    def get_reclaim_policy(self):
        ''' get reclaim policy'''
        return self.get(StorageClass.reclaim_policy_path)
