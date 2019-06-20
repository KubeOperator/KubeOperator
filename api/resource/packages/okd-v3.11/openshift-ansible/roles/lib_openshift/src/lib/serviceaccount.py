# pylint: skip-file
# flake8: noqa

class ServiceAccountConfig(object):
    '''Service account config class

       This class stores the options and returns a default service account
    '''

    # pylint: disable=too-many-arguments
    def __init__(self, sname, namespace, kubeconfig, secrets=None, image_pull_secrets=None):
        self.name = sname
        self.kubeconfig = kubeconfig
        self.namespace = namespace
        self.secrets = secrets or []
        self.image_pull_secrets = image_pull_secrets or []
        self.data = {}
        self.create_dict()

    def create_dict(self):
        ''' instantiate a properly structured volume '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'ServiceAccount'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        self.data['metadata']['namespace'] = self.namespace

        self.data['secrets'] = []
        if self.secrets:
            for sec in self.secrets:
                self.data['secrets'].append({"name": sec})

        self.data['imagePullSecrets'] = []
        if self.image_pull_secrets:
            for sec in self.image_pull_secrets:
                self.data['imagePullSecrets'].append({"name": sec})

class ServiceAccount(Yedit):
    ''' Class to wrap the oc command line tools '''
    image_pull_secrets_path = "imagePullSecrets"
    secrets_path = "secrets"

    def __init__(self, content):
        '''ServiceAccount constructor'''
        super(ServiceAccount, self).__init__(content=content)
        self._secrets = None
        self._image_pull_secrets = None

    @property
    def image_pull_secrets(self):
        ''' property for image_pull_secrets '''
        if self._image_pull_secrets is None:
            self._image_pull_secrets = self.get(ServiceAccount.image_pull_secrets_path) or []
        return self._image_pull_secrets

    @image_pull_secrets.setter
    def image_pull_secrets(self, secrets):
        ''' property for secrets '''
        self._image_pull_secrets = secrets

    @property
    def secrets(self):
        ''' property for secrets '''
        if not self._secrets:
            self._secrets = self.get(ServiceAccount.secrets_path) or []
        return self._secrets

    @secrets.setter
    def secrets(self, secrets):
        ''' property for secrets '''
        self._secrets = secrets

    def delete_secret(self, inc_secret):
        ''' remove a secret '''
        remove_idx = None
        for idx, sec in enumerate(self.secrets):
            if sec['name'] == inc_secret:
                remove_idx = idx
                break

        if remove_idx:
            del self.secrets[remove_idx]
            return True

        return False

    def delete_image_pull_secret(self, inc_secret):
        ''' remove a image_pull_secret '''
        remove_idx = None
        for idx, sec in enumerate(self.image_pull_secrets):
            if sec['name'] == inc_secret:
                remove_idx = idx
                break

        if remove_idx:
            del self.image_pull_secrets[remove_idx]
            return True

        return False

    def find_secret(self, inc_secret):
        '''find secret'''
        for secret in self.secrets:
            if secret['name'] == inc_secret:
                return secret

        return None

    def find_image_pull_secret(self, inc_secret):
        '''find secret'''
        for secret in self.image_pull_secrets:
            if secret['name'] == inc_secret:
                return secret

        return None

    def add_secret(self, inc_secret):
        '''add secret'''
        if self.secrets:
            self.secrets.append({"name": inc_secret})  # pylint: disable=no-member
        else:
            self.put(ServiceAccount.secrets_path, [{"name": inc_secret}])

    def add_image_pull_secret(self, inc_secret):
        '''add image_pull_secret'''
        if self.image_pull_secrets:
            self.image_pull_secrets.append({"name": inc_secret})  # pylint: disable=no-member
        else:
            self.put(ServiceAccount.image_pull_secrets_path, [{"name": inc_secret}])
