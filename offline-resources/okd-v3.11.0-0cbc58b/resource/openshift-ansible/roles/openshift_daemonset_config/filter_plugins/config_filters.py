#!/usr/bin/python
# -*- coding: utf-8 -*-
# pylint: disable=too-many-lines
"""
Custom filters for use in openshift-ansible
"""

from ansible import errors


def odc_join_files_from_dict(files, inc_dict):
    '''Take a list of dictionaries with name, path and insert them into
       inc_dict[name] = path
    '''
    if not isinstance(files, list):
        raise errors.AnsibleFilterError("|failed expects files param to be a list of dicts")

    if not isinstance(inc_dict, dict):
        raise errors.AnsibleFilterError("|failed expects inc_dict param to be a dict")

    for item in files:
        inc_dict[item['name']] = item['path']

    return inc_dict


class FilterModule(object):
    """ Custom ansible filter mapping """

    # pylint: disable=no-self-use, too-few-public-methods
    def filters(self):
        """ returns a mapping of filters to methods """
        return {
            "odc_join_files_from_dict": odc_join_files_from_dict,
        }
