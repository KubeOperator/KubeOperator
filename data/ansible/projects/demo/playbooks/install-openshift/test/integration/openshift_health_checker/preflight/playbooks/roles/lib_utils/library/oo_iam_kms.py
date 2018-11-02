#!/usr/bin/env python
'''
ansible module for creating AWS IAM KMS keys
'''
# vim: expandtab:tabstop=4:shiftwidth=4
#
#   AWS IAM KMS ansible module
#
#
#   Copyright 2016 Red Hat Inc.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# Jenkins environment doesn't have all the required libraries
# pylint: disable=import-error
import time
import boto3
# Ansible modules need this wildcard import
# pylint: disable=unused-wildcard-import, wildcard-import, redefined-builtin
from ansible.module_utils.basic import AnsibleModule

AWS_ALIAS_URL = "http://docs.aws.amazon.com/kms/latest/developerguide/programming-aliases.html"


class AwsIamKms(object):
    '''
    ansible module for AWS IAM KMS
    '''

    def __init__(self):
        ''' constructor '''
        self.module = None
        self.kms_client = None
        self.aliases = None

    @staticmethod
    def valid_alias_name(user_alias):
        ''' AWS KMS aliases must start with 'alias/' '''
        valid_start = 'alias/'
        if user_alias.startswith(valid_start):
            return True

        return False

    def get_all_kms_info(self):
        '''fetch all kms info and return them

        list_keys doesn't have information regarding aliases
        list_aliases doesn't have the full kms arn

        fetch both and join them on the targetKeyId
        '''
        aliases = self.kms_client.list_aliases()['Aliases']
        keys = self.kms_client.list_keys()['Keys']

        for alias in aliases:
            for key in keys:
                if 'TargetKeyId' in alias and 'KeyId' in key:
                    if alias['TargetKeyId'] == key['KeyId']:
                        alias.update(key)

        return aliases

    def get_kms_entry(self, user_alias, alias_list):
        ''' return single alias details from list of aliases '''
        for alias in alias_list:
            if user_alias == alias.get('AliasName', False):
                return alias

        msg = "Did not find alias {}".format(user_alias)
        self.module.exit_json(failed=True, results=msg)

    @staticmethod
    def exists(user_alias, alias_list):
        ''' Check if KMS alias already exists '''
        for alias in alias_list:
            if user_alias == alias.get('AliasName'):
                return True

        return False

    def main(self):
        ''' entry point for module '''

        self.module = AnsibleModule(
            argument_spec=dict(
                state=dict(default='list', choices=['list', 'present'], type='str'),
                region=dict(default=None, required=True, type='str'),
                alias=dict(default=None, type='str'),
                # description default cannot be None
                description=dict(default='', type='str'),
                aws_access_key=dict(default=None, type='str'),
                aws_secret_key=dict(default=None, type='str'),
            ),
        )

        state = self.module.params['state']
        aws_access_key = self.module.params['aws_access_key']
        aws_secret_key = self.module.params['aws_secret_key']
        if aws_access_key and aws_secret_key:
            boto3.setup_default_session(aws_access_key_id=aws_access_key,
                                        aws_secret_access_key=aws_secret_key,
                                        region_name=self.module.params['region'])
        else:
            boto3.setup_default_session(region_name=self.module.params['region'])

        self.kms_client = boto3.client('kms')

        aliases = self.get_all_kms_info()

        if state == 'list':
            if self.module.params['alias'] is not None:
                user_kms = self.get_kms_entry(self.module.params['alias'],
                                              aliases)
                self.module.exit_json(changed=False, results=user_kms,
                                      state="list")
            else:
                self.module.exit_json(changed=False, results=aliases,
                                      state="list")

        if state == 'present':

            # early sanity check to make sure the alias name conforms with
            # AWS alias name requirements
            if not self.valid_alias_name(self.module.params['alias']):
                self.module.exit_json(failed=True, changed=False,
                                      results="Alias must start with the prefix " +
                                      "'alias/'. Please see " + AWS_ALIAS_URL,
                                      state='present')

            if not self.exists(self.module.params['alias'], aliases):
                # if we didn't find it, create it
                response = self.kms_client.create_key(KeyUsage='ENCRYPT_DECRYPT',
                                                      Description=self.module.params['description'])
                kid = response['KeyMetadata']['KeyId']
                response = self.kms_client.create_alias(AliasName=self.module.params['alias'],
                                                        TargetKeyId=kid)
                # sleep for a bit so that the KMS data can be queried
                time.sleep(10)
                # get details for newly created KMS entry
                new_alias_list = self.kms_client.list_aliases()['Aliases']
                user_kms = self.get_kms_entry(self.module.params['alias'],
                                              new_alias_list)

                self.module.exit_json(changed=True, results=user_kms,
                                      state='present')

            # already exists, normally we would check whether we need to update it
            # but this module isn't written to allow changing the alias name
            # or changing whether the key is enabled/disabled
            user_kms = self.get_kms_entry(self.module.params['alias'], aliases)
            self.module.exit_json(changed=False, results=user_kms,
                                  state="present")

        self.module.exit_json(failed=True,
                              changed=False,
                              results='Unknown state passed. %s' % state,
                              state="unknown")


if __name__ == '__main__':
    AwsIamKms().main()
