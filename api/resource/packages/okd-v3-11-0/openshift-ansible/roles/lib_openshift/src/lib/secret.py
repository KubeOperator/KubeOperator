# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class SecretConfig(object):
    ''' Handle secret options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 kubeconfig,
                 secrets=None,
                 stype=None,
                 annotations=None):
        ''' constructor for handling secret options '''
        self.kubeconfig = kubeconfig
        self.name = sname
        self.type = stype
        self.namespace = namespace
        self.secrets = secrets
        self.annotations = annotations
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' assign the correct properties for a secret dict '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'Secret'
        self.data['type'] = self.type
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        self.data['metadata']['namespace'] = self.namespace
        self.data['data'] = {}
        if self.secrets:
            for key, value in self.secrets.items():
                self.data['data'][key] = value
        if self.annotations:
            self.data['metadata']['annotations'] = self.annotations

# pylint: disable=too-many-instance-attributes
class Secret(Yedit):
    ''' Class to wrap the oc command line tools '''
    secret_path = "data"
    kind = 'secret'

    def __init__(self, content):
        '''secret constructor'''
        super(Secret, self).__init__(content=content)
        self._secrets = None

    @property
    def secrets(self):
        '''secret property getter'''
        if self._secrets is None:
            self._secrets = self.get_secrets()
        return self._secrets

    @secrets.setter
    def secrets(self):
        '''secret property setter'''
        if self._secrets is None:
            self._secrets = self.get_secrets()
        return self._secrets

    def get_secrets(self):
        ''' returns all of the defined secrets '''
        return self.get(Secret.secret_path) or {}

    def add_secret(self, key, value):
        ''' add a secret '''
        if self.secrets:
            self.secrets[key] = value
        else:
            self.put(Secret.secret_path, {key: value})

        return True

    def delete_secret(self, key):
        ''' delete secret'''
        try:
            del self.secrets[key]
        except KeyError as _:
            return False

        return True

    def find_secret(self, key):
        ''' find secret'''
        rval = None
        try:
            rval = self.secrets[key]
        except KeyError as _:
            return None

        return {'key': key, 'value': rval}

    def update_secret(self, key, value):
        ''' update a secret'''
        if key in self.secrets:
            self.secrets[key] = value
        else:
            self.add_secret(key, value)

        return True
