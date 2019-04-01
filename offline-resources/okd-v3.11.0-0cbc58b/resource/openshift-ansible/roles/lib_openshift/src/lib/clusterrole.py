# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-public-methods
class ClusterRole(Yedit):
    ''' Class to model an openshift ClusterRole'''
    rules_path = "rules"

    def __init__(self, name=None, content=None):
        ''' Constructor for clusterrole '''
        if content is None:
            content = ClusterRole.builder(name).yaml_dict

        super(ClusterRole, self).__init__(content=content)

        self.__rules = Rule.parse_rules(self.get(ClusterRole.rules_path)) or []

    @property
    def rules(self):
        return self.__rules

    @rules.setter
    def rules(self, data):
        self.__rules = data
        self.put(ClusterRole.rules_path, self.__rules)

    def rule_exists(self, inc_rule):
        '''attempt to find the inc_rule in the rules list'''
        for rule in self.rules:
            if rule == inc_rule:
                return True

        return False

    def compare(self, other, verbose=False):
        '''compare function for clusterrole'''
        for rule in other.rules:
            if rule not in self.rules:
                if verbose:
                    print('Rule in other not found in self. [{}]'.format(rule))
                return False

        for rule in self.rules:
            if rule not in other.rules:
                if verbose:
                    print('Rule in self not found in other. [{}]'.format(rule))
                return False

        return True

    @staticmethod
    def builder(name='default_clusterrole', rules=None):
        '''return a clusterrole with name and/or rules'''
        if rules is None:
            rules = [{'apiGroups': [""],
                      'attributeRestrictions': None,
                      'verbs': [],
                      'resources': []}]
        content = {
            'apiVersion': 'v1',
            'kind': 'ClusterRole',
            'metadata': {'name': '{}'.format(name)},
            'rules': rules,
        }

        return ClusterRole(content=content)

