# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class SecurityContextConstraintsConfig(object):
    ''' Handle scc options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 kubeconfig,
                 options=None,
                 fs_group='MustRunAs',
                 default_add_capabilities=None,
                 groups=None,
                 priority=None,
                 required_drop_capabilities=None,
                 run_as_user='MustRunAsRange',
                 se_linux_context='MustRunAs',
                 supplemental_groups='RunAsAny',
                 users=None,
                 annotations=None):
        ''' constructor for handling scc options '''
        self.kubeconfig = kubeconfig
        self.name = sname
        self.options = options
        self.fs_group = fs_group
        self.default_add_capabilities = default_add_capabilities
        self.groups = groups
        self.priority = priority
        self.required_drop_capabilities = required_drop_capabilities
        self.run_as_user = run_as_user
        self.se_linux_context = se_linux_context
        self.supplemental_groups = supplemental_groups
        self.users = users
        self.annotations = annotations
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' assign the correct properties for a scc dict '''
        # allow options
        if self.options:
            for key, value in self.options.items():
                self.data[key] = value
        else:
            self.data['allowHostDirVolumePlugin'] = False
            self.data['allowHostIPC'] = False
            self.data['allowHostNetwork'] = False
            self.data['allowHostPID'] = False
            self.data['allowHostPorts'] = False
            self.data['allowPrivilegedContainer'] = False
            self.data['allowedCapabilities'] = None

        # version
        self.data['apiVersion'] = 'v1'
        # kind
        self.data['kind'] = 'SecurityContextConstraints'
        # defaultAddCapabilities
        self.data['defaultAddCapabilities'] = self.default_add_capabilities
        # fsGroup
        self.data['fsGroup']['type'] = self.fs_group
        # groups
        self.data['groups'] = []
        if self.groups:
            self.data['groups'] = self.groups
        # metadata
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        if self.annotations:
            for key, value in self.annotations.items():
                self.data['metadata'][key] = value
        # priority
        self.data['priority'] = self.priority
        # requiredDropCapabilities
        self.data['requiredDropCapabilities'] = self.required_drop_capabilities
        # runAsUser
        self.data['runAsUser'] = {'type': self.run_as_user}
        # seLinuxContext
        self.data['seLinuxContext'] = {'type': self.se_linux_context}
        # supplementalGroups
        self.data['supplementalGroups'] = {'type': self.supplemental_groups}
        # users
        self.data['users'] = []
        if self.users:
            self.data['users'] = self.users


# pylint: disable=too-many-instance-attributes,too-many-public-methods,no-member
class SecurityContextConstraints(Yedit):
    ''' Class to wrap the oc command line tools '''
    default_add_capabilities_path = "defaultAddCapabilities"
    fs_group_path = "fsGroup"
    groups_path = "groups"
    priority_path = "priority"
    required_drop_capabilities_path = "requiredDropCapabilities"
    run_as_user_path = "runAsUser"
    se_linux_context_path = "seLinuxContext"
    supplemental_groups_path = "supplementalGroups"
    users_path = "users"
    kind = 'SecurityContextConstraints'

    def __init__(self, content):
        '''SecurityContextConstraints constructor'''
        super(SecurityContextConstraints, self).__init__(content=content)
        self._users = None
        self._groups = None

    @property
    def users(self):
        ''' users property getter '''
        if self._users is None:
            self._users = self.get_users()
        return self._users

    @property
    def groups(self):
        ''' groups property getter '''
        if self._groups is None:
            self._groups = self.get_groups()
        return self._groups

    @users.setter
    def users(self, data):
        ''' users property setter'''
        self._users = data

    @groups.setter
    def groups(self, data):
        ''' groups property setter'''
        self._groups = data

    def get_users(self):
        '''get scc users'''
        return self.get(SecurityContextConstraints.users_path) or []

    def get_groups(self):
        '''get scc groups'''
        return self.get(SecurityContextConstraints.groups_path) or []

    def add_user(self, inc_user):
        ''' add a user '''
        if self.users:
            self.users.append(inc_user)
        else:
            self.put(SecurityContextConstraints.users_path, [inc_user])

        return True

    def add_group(self, inc_group):
        ''' add a group '''
        if self.groups:
            self.groups.append(inc_group)
        else:
            self.put(SecurityContextConstraints.groups_path, [inc_group])

        return True

    def remove_user(self, inc_user):
        ''' remove a user '''
        try:
            self.users.remove(inc_user)
        except ValueError as _:
            return False

        return True

    def remove_group(self, inc_group):
        ''' remove a group '''
        try:
            self.groups.remove(inc_group)
        except ValueError as _:
            return False

        return True

    def update_user(self, inc_user):
        ''' update a user '''
        try:
            index = self.users.index(inc_user)
        except ValueError as _:
            return self.add_user(inc_user)

        self.users[index] = inc_user

        return True

    def update_group(self, inc_group):
        ''' update a group '''
        try:
            index = self.groups.index(inc_group)
        except ValueError as _:
            return self.add_group(inc_group)

        self.groups[index] = inc_group

        return True

    def find_user(self, inc_user):
        ''' find a user '''
        index = None
        try:
            index = self.users.index(inc_user)
        except ValueError as _:
            return index

        return index

    def find_group(self, inc_group):
        ''' find a group '''
        index = None
        try:
            index = self.groups.index(inc_group)
        except ValueError as _:
            return index

        return index
