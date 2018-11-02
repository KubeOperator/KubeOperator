# pylint: skip-file
# flake8: noqa


class UserConfig(object):
    ''' Handle user options '''
    def __init__(self,
                 kubeconfig,
                 username,
                 full_name):
        ''' constructor for handling user options '''
        self.kubeconfig = kubeconfig
        self.username = username
        self.full_name = full_name

        self.data = {}
        self.create_dict()

    def create_dict(self):
        ''' return a user as a dict '''
        self.data['apiVersion'] = 'v1'
        self.data['fullName'] = self.full_name
        self.data['groups'] = None
        self.data['identities'] = None
        self.data['kind'] = 'User'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.username


# pylint: disable=too-many-instance-attributes
class User(Yedit):
    ''' Class to wrap the oc command line tools '''
    kind = 'user'

    def __init__(self, content):
        '''User constructor'''
        super(User, self).__init__(content=content)
