#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
Custom filters for use in openshift_hosted
'''


class FilterModule(object):
    ''' Custom ansible filters for use by openshift_hosted role'''

    @staticmethod
    def get_router_replicas(replicas=None, router_nodes=None):
        ''' This function will return the number of replicas
            based on the results from the defined
            openshift_hosted_router_replicas OR
            the query from oc_obj on openshift nodes with a selector OR
            default to 1

        '''
        # We always use what they've specified if they've specified a value
        if replicas is not None:
            return replicas

        replicas = 1

        # Ignore boolean expression limit of 5.
        # pylint: disable=too-many-boolean-expressions
        if (isinstance(router_nodes, dict) and
                'results' in router_nodes and
                'results' in router_nodes['results'] and
                isinstance(router_nodes['results']['results'], list) and
                len(router_nodes['results']['results']) > 0 and
                'items' in router_nodes['results']['results'][0]):

            if len(router_nodes['results']['results'][0]['items']) > 0:
                replicas = len(router_nodes['results']['results'][0]['items'])

        return replicas

    def filters(self):
        ''' returns a mapping of filters to methods '''
        return {'get_router_replicas': self.get_router_replicas}
