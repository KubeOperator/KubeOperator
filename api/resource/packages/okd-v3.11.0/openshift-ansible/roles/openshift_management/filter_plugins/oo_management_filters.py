#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
Filter methods for the management role
"""


def oo_filter_container_providers(results):
    """results - the result from posting the API calls for adding new
providers"""
    all_results = []
    for result in results:
        if 'results' in result['json']:
            # We got an OK response
            res = result['json']['results'][0]
            all_results.append("Provider '{}' - Added successfully".format(res['name']))
        elif 'error' in result['json']:
            # This was a problem
            all_results.append("Provider '{}' - Failed to add. Message: {}".format(
                result['item']['name'], result['json']['error']['message']))
    return all_results


class FilterModule(object):
    """ Custom ansible filter mapping """

    # pylint: disable=no-self-use, too-few-public-methods
    def filters(self):
        """ returns a mapping of filters to methods """
        return {
            "oo_filter_container_providers": oo_filter_container_providers,
        }
