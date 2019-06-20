#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
Custom filters for use in openshift-ansible
'''


class FilterModule(object):
    ''' Custom ansible filters '''

    @staticmethod
    def translate_volume_name(volumes, target_volume):
        '''
            This filter matches a device string /dev/sdX to /dev/xvdX
            It will then return the AWS volume ID
        '''
        for vol in volumes:
            translated_name = vol["attachment_set"]["device"].replace("/dev/sd", "/dev/xvd")
            if target_volume.startswith(translated_name):
                return vol["id"]

        return None

    def filters(self):
        ''' returns a mapping of filters to methods '''
        return {
            "translate_volume_name": self.translate_volume_name,
        }
