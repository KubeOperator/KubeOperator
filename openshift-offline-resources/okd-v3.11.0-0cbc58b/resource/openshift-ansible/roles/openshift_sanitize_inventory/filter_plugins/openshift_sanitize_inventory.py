'''
 Openshift Sanitize inventory class that provides useful filters used in Logging.
'''


import re


def vars_with_pattern(source, pattern=""):
    ''' Returns a list of variables whose name matches the given pattern '''
    if source == '':
        return list()

    var_list = list()

    var_pattern = re.compile(pattern)

    for item in source:
        if var_pattern.match(item):
            var_list.append(item)

    return var_list


# pylint: disable=too-few-public-methods
class FilterModule(object):
    ''' OpenShift Logging Filters '''

    # pylint: disable=no-self-use, too-few-public-methods
    def filters(self):
        ''' Returns the names of the filters provided by this class '''
        return {
            'vars_with_pattern': vars_with_pattern
        }
