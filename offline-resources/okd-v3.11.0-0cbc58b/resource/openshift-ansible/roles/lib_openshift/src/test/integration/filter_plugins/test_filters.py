#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
Custom filters for use in testing
'''


class FilterModule(object):
    ''' Custom filters for use in integration testing '''

    @staticmethod
    def label_dict_to_key_value_list(label_dict):
        ''' Given a dict of labels/values, return list of key: <key> value: <value> pairs

            These are only used in integration testing.
        '''

        label_list = []
        for key in label_dict:
            label_list.append({'key': key, 'value': label_dict[key]})

        return label_list

    def filters(self):
        ''' returns a mapping of filters to methods '''
        return {
            "label_dict_to_key_value_list": self.label_dict_to_key_value_list,
        }
