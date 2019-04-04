# pylint: skip-file
# flake8: noqa


class Rule(object):
    '''class to represent a clusterrole rule

    Example Rule Object's yaml:
    - apiGroups:
    - ""
    attributeRestrictions: null
    resources:
    - persistentvolumes
    verbs:
    - create
    - delete
    - deletecollection
    - get
    - list
    - patch
    - update
    - watch

    '''
    def __init__(self,
                 api_groups=None,
                 attr_restrictions=None,
                 resources=None,
                 verbs=None):
        self.__api_groups = api_groups if api_groups is not None else [""]
        self.__verbs = verbs if verbs is not None else []
        self.__resources = resources if resources is not None else []
        self.__attribute_restrictions = attr_restrictions if attr_restrictions is not None else None

    @property
    def verbs(self):
        '''property for verbs'''
        if self.__verbs is None:
            return []

        return self.__verbs

    @verbs.setter
    def verbs(self, data):
        '''setter for verbs'''
        self.__verbs = data

    @property
    def api_groups(self):
        '''property for api_groups'''
        if self.__api_groups is None:
            return []
        return self.__api_groups

    @api_groups.setter
    def api_groups(self, data):
        '''setter for api_groups'''
        self.__api_groups = data

    @property
    def resources(self):
        '''property for resources'''
        if self.__resources is None:
            return []

        return self.__resources

    @resources.setter
    def resources(self, data):
        '''setter for resources'''
        self.__resources = data

    @property
    def attribute_restrictions(self):
        '''property for attribute_restrictions'''
        return self.__attribute_restrictions

    @attribute_restrictions.setter
    def attribute_restrictions(self, data):
        '''setter for attribute_restrictions'''
        self.__attribute_restrictions = data

    def add_verb(self, inc_verb):
        '''add a verb to the verbs array'''
        self.verbs.append(inc_verb)

    def add_api_group(self, inc_apigroup):
        '''add an api_group to the api_groups array'''
        self.api_groups.append(inc_apigroup)

    def add_resource(self, inc_resource):
        '''add an resource to the resources array'''
        self.resources.append(inc_resource)

    def remove_verb(self, inc_verb):
        '''add a verb to the verbs array'''
        try:
            self.verbs.remove(inc_verb)
            return True
        except ValueError:
            pass

        return False

    def remove_api_group(self, inc_api_group):
        '''add a verb to the verbs array'''
        try:
            self.api_groups.remove(inc_api_group)
            return True
        except ValueError:
            pass

        return False

    def remove_resource(self, inc_resource):
        '''add a verb to the verbs array'''
        try:
            self.resources.remove(inc_resource)
            return True
        except ValueError:
            pass

        return False

    def __eq__(self, other):
        '''return whether rules are equal'''
        return (self.attribute_restrictions == other.attribute_restrictions and
                self.api_groups == other.api_groups and
                self.resources == other.resources and
                self.verbs == other.verbs)


    @staticmethod
    def parse_rules(inc_rules):
        '''create rules from an array'''

        results = []
        for rule in inc_rules:
            results.append(Rule(rule.get('apiGroups', ['']),
                                rule.get('attributeRestrictions', None),
                                rule.get('resources', []),
                                rule.get('verbs', [])))

        return results
