# pylint: skip-file
# flake8: noqa


class GroupConfig(object):
    ''' Handle route options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 kubeconfig):
        ''' constructor for handling group options '''
        self.kubeconfig = kubeconfig
        self.name = sname
        self.namespace = namespace
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' return a service as a dict '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'Group'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        self.data['users'] = None


# pylint: disable=too-many-instance-attributes
class Group(Yedit):
    ''' Class to wrap the oc command line tools '''
    kind = 'group'

    def __init__(self, content):
        '''Group constructor'''
        super(Group, self).__init__(content=content)
