# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class PersistentVolumeClaimConfig(object):
    ''' Handle pvc options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 kubeconfig,
                 access_modes=None,
                 vol_capacity='1G',
                 selector=None,
                 storage_class_name=None):
        ''' constructor for handling pvc options '''
        self.kubeconfig = kubeconfig
        self.name = sname
        self.namespace = namespace
        self.access_modes = access_modes
        self.vol_capacity = vol_capacity
        self.data = {}
        self.selector = selector
        self.storage_class_name = storage_class_name

        self.create_dict()

    def create_dict(self):
        ''' return a service as a dict '''
        # version
        self.data['apiVersion'] = 'v1'
        # kind
        self.data['kind'] = 'PersistentVolumeClaim'
        # metadata
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        # spec
        self.data['spec'] = {}
        self.data['spec']['accessModes'] = ['ReadWriteOnce']
        if self.access_modes:
            self.data['spec']['accessModes'] = self.access_modes
        if self.selector:
            self.data['spec']['selector'] = {'matchLabels': self.selector}

        # storage capacity
        self.data['spec']['resources'] = {}
        self.data['spec']['resources']['requests'] = {}
        self.data['spec']['resources']['requests']['storage'] = self.vol_capacity

        if self.storage_class_name:
            self.data['spec']['storageClassName'] = self.storage_class_name

# pylint: disable=too-many-instance-attributes,too-many-public-methods
class PersistentVolumeClaim(Yedit):
    ''' Class to wrap the oc command line tools '''
    access_modes_path = "spec.accessModes"
    volume_capacity_path = "spec.requests.storage"
    volume_name_path = "spec.volumeName"
    bound_path = "status.phase"
    kind = 'PersistentVolumeClaim'
    selector_path = "spec.selector.matchLabels"
    storage_class_name_path = "spec.storageClassName"

    def __init__(self, content):
        '''PersistentVolumeClaim constructor'''
        super(PersistentVolumeClaim, self).__init__(content=content)
        self._access_modes = None
        self._volume_capacity = None
        self._volume_name = None
        self._selector = None
        self._storage_class_name = None

    @property
    def storage_class_name(self):
        ''' storage_class_name property '''
        if self._storage_class_name is None:
            self._storage_class_name = self.get_storage_class_name()
        return self._storage_class_name

    @storage_class_name.setter
    def storage_class_name(self, data):
        ''' storage_class_name property setter'''
        self._storage_class_name = data

    @property
    def volume_name(self):
        ''' volume_name property '''
        if self._volume_name is None:
            self._volume_name = self.get_volume_name()
        return self._volume_name

    @volume_name.setter
    def volume_name(self, data):
        ''' volume_name property setter'''
        self._volume_name = data

    @property
    def selector(self):
        ''' selector property '''
        if self._selector is None:
            self._selector = self.get_selector()
            if not isinstance(self._selector, dict):
                self._selector = dict(self._selector)

        return self._selector

    @selector.setter
    def selector(self, data):
        ''' selector property setter'''
        if not isinstance(data, dict):
            data = dict(data)

        self._selector = data

    @property
    def access_modes(self):
        ''' access_modes property '''
        if self._access_modes is None:
            self._access_modes = self.get_access_modes()
            if not isinstance(self._access_modes, list):
                self._access_modes = list(self._access_modes)

        return self._access_modes

    @access_modes.setter
    def access_modes(self, data):
        ''' access_modes property setter'''
        if not isinstance(data, list):
            data = list(data)

        self._access_modes = data

    @property
    def volume_capacity(self):
        ''' volume_capacity property '''
        if self._volume_capacity is None:
            self._volume_capacity = self.get_volume_capacity()
        return self._volume_capacity

    @volume_capacity.setter
    def volume_capacity(self, data):
        ''' volume_capacity property setter'''
        self._volume_capacity = data

    def get_storage_class_name(self):
        '''get storage_class_name'''
        return self.get(PersistentVolumeClaim.storage_class_name_path) or []

    def get_selector(self):
        '''get selector'''
        return self.get(PersistentVolumeClaim.selector_path) or []

    def get_access_modes(self):
        '''get access_modes'''
        return self.get(PersistentVolumeClaim.access_modes_path) or []

    def get_volume_capacity(self):
        '''get volume_capacity'''
        return self.get(PersistentVolumeClaim.volume_capacity_path) or []

    def get_volume_name(self):
        '''get volume_name'''
        return self.get(PersistentVolumeClaim.volume_name_path) or []

    def is_bound(self):
        '''return whether volume is bound'''
        return self.get(PersistentVolumeClaim.bound_path) or []

    #### ADD #####
    def add_access_mode(self, inc_mode):
        ''' add an access_mode'''
        if self.access_modes:
            self.access_modes.append(inc_mode)
        else:
            self.put(PersistentVolumeClaim.access_modes_path, [inc_mode])

        return True

    #### /ADD #####

    #### Remove #####
    def remove_access_mode(self, inc_mode):
        ''' remove an access_mode'''
        try:
            self.access_modes.remove(inc_mode)
        except ValueError as _:
            return False

        return True

    #### /REMOVE #####

    #### UPDATE #####
    def update_access_mode(self, inc_mode):
        ''' update an access_mode'''
        try:
            index = self.access_modes.index(inc_mode)
        except ValueError as _:
            return self.add_access_mode(inc_mode)

        self.access_modes[index] = inc_mode

        return True

    #### /UPDATE #####

    #### FIND ####
    def find_access_mode(self, inc_mode):
        ''' find a user '''
        index = None
        try:
            index = self.access_modes.index(inc_mode)
        except ValueError as _:
            return index

        return index
