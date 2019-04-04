#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
Custom filters for use in openshift-ansible
"""


# Disabling too-many-public-methods, since filter methods are necessarily
# public
# pylint: disable=too-many-public-methods
class FilterModule(object):
    """ Custom ansible filters """

    @staticmethod
    def oo_cert_expiry_results_to_json(hostvars, play_hosts):
        """Takes results (`hostvars`) from the openshift_cert_expiry role
check and serializes them into proper machine-readable JSON
output. This filter parameter **MUST** be the playbook `hostvars`
variable. The `play_hosts` parameter is so we know what to loop over
when we're extrating the values.

Returns:

Results are collected into two top-level keys under the `json_results`
dict:

* `json_results.data` [dict] - Each individual host check result, keys are hostnames
* `json_results.summary` [dict] - Summary of number of `warning` and `expired`
certificates

Example playbook usage:

  - name: Generate expiration results JSON
    run_once: yes
    delegate_to: localhost
    when: openshift_certificate_expiry_save_json_results|bool
    copy:
      content: "{{ hostvars|oo_cert_expiry_results_to_json() }}"
      dest: "{{ openshift_certificate_expiry_json_results_path }}"

        """
        json_result = {
            'data': {},
            'summary': {},
        }

        for host in play_hosts:
            json_result['data'][host] = hostvars[host]['check_results']['check_results']

        total_warnings = sum([hostvars[h]['check_results']['summary']['warning'] for h in play_hosts])
        total_expired = sum([hostvars[h]['check_results']['summary']['expired'] for h in play_hosts])
        total_ok = sum([hostvars[h]['check_results']['summary']['ok'] for h in play_hosts])
        total_total = sum([hostvars[h]['check_results']['summary']['total'] for h in play_hosts])

        json_result['summary']['warning'] = total_warnings
        json_result['summary']['expired'] = total_expired
        json_result['summary']['ok'] = total_ok
        json_result['summary']['total'] = total_total

        return json_result

    def filters(self):
        """ returns a mapping of filters to methods """
        return {
            "oo_cert_expiry_results_to_json": self.oo_cert_expiry_results_to_json,
        }
