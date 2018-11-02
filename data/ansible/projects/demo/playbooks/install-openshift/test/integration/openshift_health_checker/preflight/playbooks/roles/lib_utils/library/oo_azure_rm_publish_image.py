#!/usr/bin/env python
# pylint: disable=missing-docstring
# Copyright 2018 Red Hat, Inc. and/or its affiliates
# and other contributors as indicated by the @author tags.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

from __future__ import print_function  # noqa: F401
# import httplib
import json
import os
import time
import requests

from ansible.module_utils.basic import AnsibleModule


class AzurePublisherException(Exception):
    '''Exception class for AzurePublisher'''
    pass


class AzurePublisher(object):
    '''Python class to represent the Azure Publishing portal https://cloudpartner.azure.com'''

    # pylint: disable=too-many-arguments
    def __init__(self,
                 publisher_id,
                 client_info,
                 ssl_verify=True,
                 api_version='2017-10-31',
                 debug=False):
        '''
          :publisher_id: string of the publisher id
          :client_info: a dict containing the client_id, client_secret to get an access_token
        '''
        self._azure_server = 'https://cloudpartner.azure.com/api/publishers/{}'.format(publisher_id)
        self.client_info = client_info
        self.ssl_verify = ssl_verify
        self.api_version = 'api-version={}'.format(api_version)
        self.debug = debug
        # if self.debug:
        # import httplib
        # httplib.HTTPSConnection.debuglevel = 1
        # httplib.HTTPConnection.debuglevel = 1

        self._access_token = None

    @property
    def server(self):
        '''property for  server url'''
        return self._azure_server

    @property
    def token(self):
        '''property for the access_token
            curl -d \
            'client_id=<id>&client_secret=<sec>&grant_type=client_credentials&resource=https://cloudpartner.azure.com' \
            https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47/oauth2/token
        '''
        if self._access_token is None:
            url = 'https://login.microsoftonline.com/{}/oauth2/token'.format(self.client_info['tenant_id'])
            data = {
                'client_id': {self.client_info['client_id']},
                'client_secret': self.client_info['client_secret'],
                'grant_type': 'client_credentials',
                'resource': 'https://cloudpartner.azure.com'
            }

            results = AzurePublisher.request('POST', url, data, {})
            jres = results.json()
            self._access_token = jres['access_token']

        return self._access_token

    def get_offers(self, offer=None, version=None, slot='preview'):
        ''' fetch all offers by publisherid '''
        url = '/offers'

        if offer is not None:
            url += '/{}'.format(offer)
            if version is not None:
                url += '/versions/{}'.format(version)
            if slot == 'preview':
                url += '/slot/{}'.format(slot)

        url += '?{}'.format(self.api_version)

        return self.prepare_action(url)

    def get_operations(self, offer, operation=None, status=None):
        ''' create or modify an offer '''
        url = '/offers/{0}/submissions'.format(offer)

        if operation is not None:
            url += '/operations/{0}'.format(operation)

        if not url.endswith('/'):
            url += '/'

        url += '?{0}'.format(self.api_version)

        if status is not None:
            url += '&status={0}'.format(status)

        return self.prepare_action(url, 'GET')

    def cancel_operation(self, offer):
        ''' create or modify an offer '''
        url = '/offers/{0}/cancel?{1}'.format(offer, self.api_version)

        return self.prepare_action(url, 'POST')

    def publish(self, offer, emails):
        ''' publish an offer '''
        url = '/offers/{0}/publish?{1}'.format(offer, self.api_version)

        data = {
            'metadata': {
                'notification-emails': ','.join(emails),
            }
        }

        return self.prepare_action(url, 'POST', data=data)

    def go_live(self, offer):
        ''' create or modify an offer '''
        url = '/offers/{0}/golive?{1}'.format(offer, self.api_version)

        return self.prepare_action(url, 'POST')

    def create_or_modify_offer(self, offer, data=None, modify=False):
        ''' create or modify an offer '''
        url = '/offers/{0}?{1}'.format(offer, self.api_version)

        headers = None

        if modify:
            headers = {
                'If-Match': '*',
            }

        return self.prepare_action(url, 'PUT', data=data, add_headers=headers)

    def prepare_action(self, url, action='GET', data=None, add_headers=None):
        '''perform the http request

           :action: string of either GET|POST
        '''
        headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': 'Bearer {}'.format(self.token)
        }

        if add_headers is not None:
            headers.update(add_headers)

        if data is None:
            data = ''
        else:
            data = json.dumps(data)

        return AzurePublisher.request(action.upper(), self.server + url, data, headers)

    def cancel_and_wait_for_operation(self, params):
        '''cancel the current publish operation and wait for operation to complete'''

        # cancel the publish operation
        self.cancel_operation(offer=params['offer'])

        # we need to wait here for 'submissionState' to move to 'canceled'
        while True:
            # fetch operations
            ops = self.get_operations(params['offer'])
            if self.debug:
                print(ops.json())
            if ops.json()[0]['submissionState'] == 'canceled':
                break

            time.sleep(5)

        return ops

    def manage_offer(self, params):
        ''' handle creating or modifying offers'''
        # fetch the offer to verify it exists:
        results = self.get_offers(offer=params['offer'])

        if results.status_code == 200 and params['force']:
            return self.create_or_modify_offer(offer=params['offer'], data=params['offer_data'], modify=True)

        return self.create_or_modify_offer(offer=params['offer'], data=params['offer_data'])

    @staticmethod
    def request(action, url, data=None, headers=None, ssl_verify=True):
        req = requests.Request(action.upper(), url, data=data, headers=headers)

        session = requests.Session()
        req_prep = session.prepare_request(req)
        response = session.send(req_prep, verify=ssl_verify)

        return response

    @staticmethod
    def run_ansible(params):
        '''perform the ansible operations'''
        client_info = {
            'tenant_id': params['tenant_id'],
            'client_id': params['client_id'],
            'client_secret': params['client_secret']}

        apc = AzurePublisher(params['publisher'],
                             client_info,
                             debug=params['debug'])

        if params['state'] == 'offer':
            results = apc.manage_offer(params)
        elif params['state'] == 'publish':
            results = apc.publish(offer=params['offer'], emails=params['emails'])
            results.json = lambda: ''
        elif params['state'] == 'cancel_op':
            results = apc.cancel_and_wait_for_operation(params)
        elif params['state'] == 'go_live':
            results = apc.go_live(offer=params['offer'])
        else:
            raise AzurePublisherException('Unsupported query type: {}'.format(params['state']))

        changed = False

        if results.status_code in [200, 201, 202]:
            changed = True

        return {'data': results.json(), 'changed': changed, 'status_code': results.status_code}


# pylint: disable=too-many-branches
def main():
    ''' ansible oc module for secrets '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='offer', choices=['offer', 'cancel_op', 'go_live', 'publish']),
            force=dict(default=False, type='bool'),
            publisher=dict(default='redhat', type='str'),
            debug=dict(default=False, type='bool'),
            tenant_id=dict(default=os.environ.get('AZURE_TENANT_ID'), type='str'),
            client_id=dict(default=os.environ.get('AZURE_CLIENT_ID'), type='str'),
            client_secret=dict(default=os.environ.get('AZURE_CLIENT_SECRET'), type='str'),
            offer=dict(default=None, type='str'),
            offer_data=dict(default=None, type='dict'),
            emails=dict(default=None, type='list'),
        ),

        required_if=[
            ["state", "offer", ["offer_data"]],
        ],
    )

    # Verify we recieved either a valid key or edits with valid keys when receiving a src file.
    # A valid key being not None or not ''.
    if (module.params['tenant_id'] is None or module.params['client_id'] is None or
            module.params['client_secret'] is None):
        return module.fail_json(**{'failed': True,
                                   'msg': 'Please specify tenant_id, client_id, and client_secret'})

    rval = AzurePublisher.run_ansible(module.params)

    if int(rval['status_code']) >= 300:
        rval['msg'] = 'Failed. status_code {}'.format(rval['status_code'])
        return module.fail_json(**rval)

    return module.exit_json(**rval)


if __name__ == '__main__':
    main()
