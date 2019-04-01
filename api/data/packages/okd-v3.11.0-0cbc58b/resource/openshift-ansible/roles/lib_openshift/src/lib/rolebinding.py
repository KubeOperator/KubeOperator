# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class RoleBindingConfig(object):
    ''' Handle rolebinding config '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 name,
                 namespace,
                 kubeconfig,
                 group_names=None,
                 role_ref=None,
                 subjects=None,
                 usernames=None):
        ''' constructor for handling rolebinding options '''
        self.kubeconfig = kubeconfig
        self.name = name
        self.namespace = namespace
        self.group_names = group_names
        self.role_ref = role_ref
        self.subjects = subjects
        self.usernames = usernames
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' create a default rolebinding as a dict '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'RoleBinding'
        self.data['groupNames'] = self.group_names
        self.data['metadata']['name'] = self.name
        self.data['metadata']['namespace'] = self.namespace

        self.data['roleRef'] = self.role_ref
        self.data['subjects'] = self.subjects
        self.data['userNames'] = self.usernames


# pylint: disable=too-many-instance-attributes,too-many-public-methods
class RoleBinding(Yedit):
    ''' Class to model a rolebinding openshift object'''
    group_names_path = "groupNames"
    role_ref_path = "roleRef"
    subjects_path = "subjects"
    user_names_path = "userNames"

    kind = 'RoleBinding'

    def __init__(self, content):
        '''RoleBinding constructor'''
        super(RoleBinding, self).__init__(content=content)
        self._subjects = None
        self._role_ref = None
        self._group_names = None
        self._user_names = None

    @property
    def subjects(self):
        ''' subjects property '''
        if self._subjects is None:
            self._subjects = self.get_subjects()
        return self._subjects

    @subjects.setter
    def subjects(self, data):
        ''' subjects property setter'''
        self._subjects = data

    @property
    def role_ref(self):
        ''' role_ref property '''
        if self._role_ref is None:
            self._role_ref = self.get_role_ref()
        return self._role_ref

    @role_ref.setter
    def role_ref(self, data):
        ''' role_ref property setter'''
        self._role_ref = data

    @property
    def group_names(self):
        ''' group_names property '''
        if self._group_names is None:
            self._group_names = self.get_group_names()
        return self._group_names

    @group_names.setter
    def group_names(self, data):
        ''' group_names property setter'''
        self._group_names = data

    @property
    def user_names(self):
        ''' user_names property '''
        if self._user_names is None:
            self._user_names = self.get_user_names()
        return self._user_names

    @user_names.setter
    def user_names(self, data):
        ''' user_names property setter'''
        self._user_names = data

    def get_group_names(self):
        ''' return groupNames '''
        return self.get(RoleBinding.group_names_path) or []

    def get_user_names(self):
        ''' return usernames '''
        return self.get(RoleBinding.user_names_path) or []

    def get_role_ref(self):
        ''' return role_ref '''
        return self.get(RoleBinding.role_ref_path) or {}

    def get_subjects(self):
        ''' return subjects '''
        return self.get(RoleBinding.subjects_path) or []

    #### ADD #####
    def add_subject(self, inc_subject):
        ''' add a subject '''
        if self.subjects:
            # pylint: disable=no-member
            self.subjects.append(inc_subject)
        else:
            self.put(RoleBinding.subjects_path, [inc_subject])

        return True

    def add_role_ref(self, inc_role_ref):
        ''' add a role_ref '''
        if not self.role_ref:
            self.put(RoleBinding.role_ref_path, {"name": inc_role_ref})
            return True

        return False

    def add_group_names(self, inc_group_names):
        ''' add a group_names '''
        if self.group_names:
            # pylint: disable=no-member
            self.group_names.append(inc_group_names)
        else:
            self.put(RoleBinding.group_names_path, [inc_group_names])

        return True

    def add_user_name(self, inc_user_name):
        ''' add a username '''
        if self.user_names:
            # pylint: disable=no-member
            self.user_names.append(inc_user_name)
        else:
            self.put(RoleBinding.user_names_path, [inc_user_name])

        return True

    #### /ADD #####

    #### Remove #####
    def remove_subject(self, inc_subject):
        ''' remove a subject '''
        try:
            # pylint: disable=no-member
            self.subjects.remove(inc_subject)
        except ValueError as _:
            return False

        return True

    def remove_role_ref(self, inc_role_ref):
        ''' remove a role_ref '''
        if self.role_ref and self.role_ref['name'] == inc_role_ref:
            del self.role_ref['name']
            return True

        return False

    def remove_group_name(self, inc_group_name):
        ''' remove a groupname '''
        try:
            # pylint: disable=no-member
            self.group_names.remove(inc_group_name)
        except ValueError as _:
            return False

        return True

    def remove_user_name(self, inc_user_name):
        ''' remove a username '''
        try:
            # pylint: disable=no-member
            self.user_names.remove(inc_user_name)
        except ValueError as _:
            return False

        return True

    #### /REMOVE #####

    #### UPDATE #####
    def update_subject(self, inc_subject):
        ''' update a subject '''
        try:
            # pylint: disable=no-member
            index = self.subjects.index(inc_subject)
        except ValueError as _:
            return self.add_subject(inc_subject)

        self.subjects[index] = inc_subject

        return True

    def update_group_name(self, inc_group_name):
        ''' update a groupname '''
        try:
            # pylint: disable=no-member
            index = self.group_names.index(inc_group_name)
        except ValueError as _:
            return self.add_group_names(inc_group_name)

        self.group_names[index] = inc_group_name

        return True

    def update_user_name(self, inc_user_name):
        ''' update a username '''
        try:
            # pylint: disable=no-member
            index = self.user_names.index(inc_user_name)
        except ValueError as _:
            return self.add_user_name(inc_user_name)

        self.user_names[index] = inc_user_name

        return True

    def update_role_ref(self, inc_role_ref):
        ''' update a role_ref '''
        self.role_ref['name'] = inc_role_ref

        return True

    #### /UPDATE #####

    #### FIND ####
    def find_subject(self, inc_subject):
        ''' find a subject '''
        index = None
        try:
            # pylint: disable=no-member
            index = self.subjects.index(inc_subject)
        except ValueError as _:
            return index

        return index

    def find_group_name(self, inc_group_name):
        ''' find a group_name '''
        index = None
        try:
            # pylint: disable=no-member
            index = self.group_names.index(inc_group_name)
        except ValueError as _:
            return index

        return index

    def find_user_name(self, inc_user_name):
        ''' find a user_name '''
        index = None
        try:
            # pylint: disable=no-member
            index = self.user_names.index(inc_user_name)
        except ValueError as _:
            return index

        return index

    def find_role_ref(self, inc_role_ref):
        ''' find a user_name '''
        if self.role_ref and self.role_ref['name'] == inc_role_ref['name']:
            return self.role_ref

        return None
