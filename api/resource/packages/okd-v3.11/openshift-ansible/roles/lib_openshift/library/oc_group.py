#!/usr/bin/env python
# pylint: disable=missing-docstring
# flake8: noqa: T001
#     ___ ___ _  _ ___ ___    _ _____ ___ ___
#    / __| __| \| | __| _ \  /_\_   _| __|   \
#   | (_ | _|| .` | _||   / / _ \| | | _|| |) |
#    \___|___|_|\_|___|_|_\/_/_\_\_|_|___|___/_ _____
#   |   \ / _ \  | \| |/ _ \_   _| | __|   \_ _|_   _|
#   | |) | (_) | | .` | (_) || |   | _|| |) | |  | |
#   |___/ \___/  |_|\_|\___/ |_|   |___|___/___| |_|
#
# Copyright 2016 Red Hat, Inc. and/or its affiliates
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

# -*- -*- -*- Begin included fragment: lib/import.py -*- -*- -*-
'''
   OpenShiftCLI class that wraps the oc commands in a subprocess
'''
# pylint: disable=too-many-lines

from __future__ import print_function
import atexit
import copy
import fcntl
import json
import time
import os
import re
import shutil
import subprocess
import tempfile
# pylint: disable=import-error
try:
    import ruamel.yaml as yaml
except ImportError:
    import yaml

from ansible.module_utils.basic import AnsibleModule

# -*- -*- -*- End included fragment: lib/import.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: doc/group -*- -*- -*-

DOCUMENTATION = '''
---
module: oc_group
short_description: Modify, and idempotently manage openshift groups.
description:
  - Modify openshift groups programmatically.
options:
  state:
    description:
    - Supported states, present, absent, list
    - present - will ensure object is created or updated to the value specified
    - list - will return a group
    - absent - will remove the group
    required: False
    default: present
    choices: ["present", 'absent', 'list']
    aliases: []
  kubeconfig:
    description:
    - The path for the kubeconfig file to use for authentication
    required: false
    default: /etc/origin/master/admin.kubeconfig
    aliases: []
  debug:
    description:
    - Turn on debug output.
    required: false
    default: False
    aliases: []
  name:
    description:
    - Name of the object that is being queried.
    required: false
    default: None
    aliases: []
  namespace:
    description:
    - The namespace where the object lives.
    required: false
    default: str
    aliases: []
author:
- "Joel Diaz <jdiaz@redhat.com>"
extends_documentation_fragment: []
'''

EXAMPLES = '''
- name: create group
  oc_group:
    state: present
    name: acme_org
  register: group_out
'''

# -*- -*- -*- End included fragment: doc/group -*- -*- -*-

# -*- -*- -*- Begin included fragment: ../../lib_utils/src/class/yedit.py -*- -*- -*-


class YeditException(Exception):  # pragma: no cover
    ''' Exception class for Yedit '''
    pass


# pylint: disable=too-many-public-methods,too-many-instance-attributes
class Yedit(object):  # pragma: no cover
    ''' Class to modify yaml files '''
    re_valid_key = r"(((\[-?\d+\])|([0-9a-zA-Z%s/_-]+)).?)+$"
    re_key = r"(?:\[(-?\d+)\])|([0-9a-zA-Z{}/_-]+)"
    com_sep = set(['.', '#', '|', ':'])

    # pylint: disable=too-many-arguments
    def __init__(self,
                 filename=None,
                 content=None,
                 content_type='yaml',
                 separator='.',
                 backup_ext=None,
                 backup=False):
        self.content = content
        self._separator = separator
        self.filename = filename
        self.__yaml_dict = content
        self.content_type = content_type
        self.backup = backup
        if backup_ext is None:
            self.backup_ext = ".{}".format(time.strftime("%Y%m%dT%H%M%S"))
        else:
            self.backup_ext = backup_ext

        self.load(content_type=self.content_type)
        if self.__yaml_dict is None:
            self.__yaml_dict = {}

    @property
    def separator(self):
        ''' getter method for separator '''
        return self._separator

    @separator.setter
    def separator(self, inc_sep):
        ''' setter method for separator '''
        self._separator = inc_sep

    @property
    def yaml_dict(self):
        ''' getter method for yaml_dict '''
        return self.__yaml_dict

    @yaml_dict.setter
    def yaml_dict(self, value):
        ''' setter method for yaml_dict '''
        self.__yaml_dict = value

    @staticmethod
    def parse_key(key, sep='.'):
        '''parse the key allowing the appropriate separator'''
        common_separators = list(Yedit.com_sep - set([sep]))
        return re.findall(Yedit.re_key.format(''.join(common_separators)), key)

    @staticmethod
    def valid_key(key, sep='.'):
        '''validate the incoming key'''
        common_separators = list(Yedit.com_sep - set([sep]))
        if not re.match(Yedit.re_valid_key.format(''.join(common_separators)), key):
            return False

        return True

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def remove_entry(data, key, index=None, value=None, sep='.'):
        ''' remove data at location key '''
        if key == '' and isinstance(data, dict):
            if value is not None:
                data.pop(value)
            elif index is not None:
                raise YeditException("remove_entry for a dictionary does not have an index {}".format(index))
            else:
                data.clear()

            return True

        elif key == '' and isinstance(data, list):
            ind = None
            if value is not None:
                try:
                    ind = data.index(value)
                except ValueError:
                    return False
            elif index is not None:
                ind = index
            else:
                del data[:]

            if ind is not None:
                data.pop(ind)

            return True

        if not (key and Yedit.valid_key(key, sep)) and \
           isinstance(data, (list, dict)):
            return None

        key_indexes = Yedit.parse_key(key, sep)
        for arr_ind, dict_key in key_indexes[:-1]:
            if dict_key and isinstance(data, dict):
                data = data.get(dict_key)
            elif (arr_ind and isinstance(data, list) and
                  int(arr_ind) <= len(data) - 1):
                data = data[int(arr_ind)]
            else:
                return None

        # process last index for remove
        # expected list entry
        if key_indexes[-1][0]:
            if isinstance(data, list) and int(key_indexes[-1][0]) <= len(data) - 1:  # noqa: E501
                del data[int(key_indexes[-1][0])]
                return True

        # expected dict entry
        elif key_indexes[-1][1]:
            if isinstance(data, dict):
                del data[key_indexes[-1][1]]
                return True

    @staticmethod
    def add_entry(data, key, item=None, sep='.'):
        ''' Get an item from a dictionary with key notation a.b.c
            d = {'a': {'b': 'c'}}}
            key = a#b
            return c
        '''
        if key == '':
            pass
        elif (not (key and Yedit.valid_key(key, sep)) and
              isinstance(data, (list, dict))):
            return None

        key_indexes = Yedit.parse_key(key, sep)
        for arr_ind, dict_key in key_indexes[:-1]:
            if dict_key:
                if isinstance(data, dict) and dict_key in data and data[dict_key]:  # noqa: E501
                    data = data[dict_key]
                    continue

                elif data and not isinstance(data, dict):
                    raise YeditException("Unexpected item type found while going through key " +
                                         "path: {} (at key: {})".format(key, dict_key))

                data[dict_key] = {}
                data = data[dict_key]

            elif (arr_ind and isinstance(data, list) and
                  int(arr_ind) <= len(data) - 1):
                data = data[int(arr_ind)]
            else:
                raise YeditException("Unexpected item type found while going through key path: {}".format(key))

        if key == '':
            data = item

        # process last index for add
        # expected list entry
        elif key_indexes[-1][0] and isinstance(data, list) and int(key_indexes[-1][0]) <= len(data) - 1:  # noqa: E501
            data[int(key_indexes[-1][0])] = item

        # expected dict entry
        elif key_indexes[-1][1] and isinstance(data, dict):
            data[key_indexes[-1][1]] = item

        # didn't add/update to an existing list, nor add/update key to a dict
        # so we must have been provided some syntax like a.b.c[<int>] = "data" for a
        # non-existent array
        else:
            raise YeditException("Error adding to object at path: {}".format(key))

        return data

    @staticmethod
    def get_entry(data, key, sep='.'):
        ''' Get an item from a dictionary with key notation a.b.c
            d = {'a': {'b': 'c'}}}
            key = a.b
            return c
        '''
        if key == '':
            pass
        elif (not (key and Yedit.valid_key(key, sep)) and
              isinstance(data, (list, dict))):
            return None

        key_indexes = Yedit.parse_key(key, sep)
        for arr_ind, dict_key in key_indexes:
            if dict_key and isinstance(data, dict):
                data = data.get(dict_key)
            elif (arr_ind and isinstance(data, list) and
                  int(arr_ind) <= len(data) - 1):
                data = data[int(arr_ind)]
            else:
                return None

        return data

    @staticmethod
    def _write(filename, contents):
        ''' Actually write the file contents to disk. This helps with mocking. '''

        tmp_filename = filename + '.yedit'

        with open(tmp_filename, 'w') as yfd:
            fcntl.flock(yfd, fcntl.LOCK_EX | fcntl.LOCK_NB)
            yfd.write(contents)
            fcntl.flock(yfd, fcntl.LOCK_UN)

        os.rename(tmp_filename, filename)

    def write(self):
        ''' write to file '''
        if not self.filename:
            raise YeditException('Please specify a filename.')

        if self.backup and self.file_exists():
            shutil.copy(self.filename, '{}{}'.format(self.filename, self.backup_ext))

        # Try to set format attributes if supported
        try:
            self.yaml_dict.fa.set_block_style()
        except AttributeError:
            pass

        # Try to use RoundTripDumper if supported.
        if self.content_type == 'yaml':
            try:
                Yedit._write(self.filename, yaml.dump(self.yaml_dict, Dumper=yaml.RoundTripDumper))
            except AttributeError:
                Yedit._write(self.filename, yaml.safe_dump(self.yaml_dict, default_flow_style=False))
        elif self.content_type == 'json':
            Yedit._write(self.filename, json.dumps(self.yaml_dict, indent=4, sort_keys=True))
        else:
            raise YeditException('Unsupported content_type: {}.'.format(self.content_type) +
                                 'Please specify a content_type of yaml or json.')

        return (True, self.yaml_dict)

    def read(self):
        ''' read from file '''
        # check if it exists
        if self.filename is None or not self.file_exists():
            return None

        contents = None
        with open(self.filename) as yfd:
            contents = yfd.read()

        return contents

    def file_exists(self):
        ''' return whether file exists '''
        if os.path.exists(self.filename):
            return True

        return False

    def load(self, content_type='yaml'):
        ''' return yaml file '''
        contents = self.read()

        if not contents and not self.content:
            return None

        if self.content:
            if isinstance(self.content, dict):
                self.yaml_dict = self.content
                return self.yaml_dict
            elif isinstance(self.content, str):
                contents = self.content

        # check if it is yaml
        try:
            if content_type == 'yaml' and contents:
                # Try to set format attributes if supported
                try:
                    self.yaml_dict.fa.set_block_style()
                except AttributeError:
                    pass

                # Try to use RoundTripLoader if supported.
                try:
                    self.yaml_dict = yaml.load(contents, yaml.RoundTripLoader)
                except AttributeError:
                    self.yaml_dict = yaml.safe_load(contents)

                # Try to set format attributes if supported
                try:
                    self.yaml_dict.fa.set_block_style()
                except AttributeError:
                    pass

            elif content_type == 'json' and contents:
                self.yaml_dict = json.loads(contents)
        except yaml.YAMLError as err:
            # Error loading yaml or json
            raise YeditException('Problem with loading yaml file. {}'.format(err))

        return self.yaml_dict

    def get(self, key):
        ''' get a specified key'''
        try:
            entry = Yedit.get_entry(self.yaml_dict, key, self.separator)
        except KeyError:
            entry = None

        return entry

    def pop(self, path, key_or_item):
        ''' remove a key, value pair from a dict or an item for a list'''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if entry is None:
            return (False, self.yaml_dict)

        if isinstance(entry, dict):
            # AUDIT:maybe-no-member makes sense due to fuzzy types
            # pylint: disable=maybe-no-member
            if key_or_item in entry:
                entry.pop(key_or_item)
                return (True, self.yaml_dict)
            return (False, self.yaml_dict)

        elif isinstance(entry, list):
            # AUDIT:maybe-no-member makes sense due to fuzzy types
            # pylint: disable=maybe-no-member
            ind = None
            try:
                ind = entry.index(key_or_item)
            except ValueError:
                return (False, self.yaml_dict)

            entry.pop(ind)
            return (True, self.yaml_dict)

        return (False, self.yaml_dict)

    def delete(self, path, index=None, value=None):
        ''' remove path from a dict'''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if entry is None:
            return (False, self.yaml_dict)

        result = Yedit.remove_entry(self.yaml_dict, path, index, value, self.separator)
        if not result:
            return (False, self.yaml_dict)

        return (True, self.yaml_dict)

    def exists(self, path, value):
        ''' check if value exists at path'''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if isinstance(entry, list):
            if value in entry:
                return True
            return False

        elif isinstance(entry, dict):
            if isinstance(value, dict):
                rval = False
                for key, val in value.items():
                    if entry[key] != val:
                        rval = False
                        break
                else:
                    rval = True
                return rval

            return value in entry

        return entry == value

    def append(self, path, value):
        '''append value to a list'''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if entry is None:
            self.put(path, [])
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        if not isinstance(entry, list):
            return (False, self.yaml_dict)

        # AUDIT:maybe-no-member makes sense due to loading data from
        # a serialized format.
        # pylint: disable=maybe-no-member
        entry.append(value)
        return (True, self.yaml_dict)

    # pylint: disable=too-many-arguments
    def update(self, path, value, index=None, curr_value=None):
        ''' put path, value into a dict '''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if isinstance(entry, dict):
            # AUDIT:maybe-no-member makes sense due to fuzzy types
            # pylint: disable=maybe-no-member
            if not isinstance(value, dict):
                raise YeditException('Cannot replace key, value entry in dict with non-dict type. ' +
                                     'value=[{}] type=[{}]'.format(value, type(value)))

            entry.update(value)
            return (True, self.yaml_dict)

        elif isinstance(entry, list):
            # AUDIT:maybe-no-member makes sense due to fuzzy types
            # pylint: disable=maybe-no-member
            ind = None
            if curr_value:
                try:
                    ind = entry.index(curr_value)
                except ValueError:
                    return (False, self.yaml_dict)

            elif index is not None:
                ind = index

            if ind is not None and entry[ind] != value:
                entry[ind] = value
                return (True, self.yaml_dict)

            # see if it exists in the list
            try:
                ind = entry.index(value)
            except ValueError:
                # doesn't exist, append it
                entry.append(value)
                return (True, self.yaml_dict)

            # already exists, return
            if ind is not None:
                return (False, self.yaml_dict)
        return (False, self.yaml_dict)

    def put(self, path, value):
        ''' put path, value into a dict '''
        try:
            entry = Yedit.get_entry(self.yaml_dict, path, self.separator)
        except KeyError:
            entry = None

        if entry == value:
            return (False, self.yaml_dict)

        # deepcopy didn't work
        # Try to use ruamel.yaml and fallback to pyyaml
        try:
            tmp_copy = yaml.load(yaml.round_trip_dump(self.yaml_dict,
                                                      default_flow_style=False),
                                 yaml.RoundTripLoader)
        except AttributeError:
            tmp_copy = copy.deepcopy(self.yaml_dict)

        # set the format attributes if available
        try:
            tmp_copy.fa.set_block_style()
        except AttributeError:
            pass

        result = Yedit.add_entry(tmp_copy, path, value, self.separator)
        if result is None:
            return (False, self.yaml_dict)

        # When path equals "" it is a special case.
        # "" refers to the root of the document
        # Only update the root path (entire document) when its a list or dict
        if path == '':
            if isinstance(result, list) or isinstance(result, dict):
                self.yaml_dict = result
                return (True, self.yaml_dict)

            return (False, self.yaml_dict)

        self.yaml_dict = tmp_copy

        return (True, self.yaml_dict)

    def create(self, path, value):
        ''' create a yaml file '''
        if not self.file_exists():
            # deepcopy didn't work
            # Try to use ruamel.yaml and fallback to pyyaml
            try:
                tmp_copy = yaml.load(yaml.round_trip_dump(self.yaml_dict,
                                                          default_flow_style=False),
                                     yaml.RoundTripLoader)
            except AttributeError:
                tmp_copy = copy.deepcopy(self.yaml_dict)

            # set the format attributes if available
            try:
                tmp_copy.fa.set_block_style()
            except AttributeError:
                pass

            result = Yedit.add_entry(tmp_copy, path, value, self.separator)
            if result is not None:
                self.yaml_dict = tmp_copy
                return (True, self.yaml_dict)

        return (False, self.yaml_dict)

    @staticmethod
    def get_curr_value(invalue, val_type):
        '''return the current value'''
        if invalue is None:
            return None

        curr_value = invalue
        if val_type == 'yaml':
            curr_value = yaml.safe_load(str(invalue))
        elif val_type == 'json':
            curr_value = json.loads(invalue)

        return curr_value

    @staticmethod
    def parse_value(inc_value, vtype=''):
        '''determine value type passed'''
        true_bools = ['y', 'Y', 'yes', 'Yes', 'YES', 'true', 'True', 'TRUE',
                      'on', 'On', 'ON', ]
        false_bools = ['n', 'N', 'no', 'No', 'NO', 'false', 'False', 'FALSE',
                       'off', 'Off', 'OFF']

        # It came in as a string but you didn't specify value_type as string
        # we will convert to bool if it matches any of the above cases
        if isinstance(inc_value, str) and 'bool' in vtype:
            if inc_value not in true_bools and inc_value not in false_bools:
                raise YeditException('Not a boolean type. str=[{}] vtype=[{}]'.format(inc_value, vtype))
        elif isinstance(inc_value, bool) and 'str' in vtype:
            inc_value = str(inc_value)

        # There is a special case where '' will turn into None after yaml loading it so skip
        if isinstance(inc_value, str) and inc_value == '':
            pass
        # If vtype is not str then go ahead and attempt to yaml load it.
        elif isinstance(inc_value, str) and 'str' not in vtype:
            try:
                inc_value = yaml.safe_load(inc_value)
            except Exception:
                raise YeditException('Could not determine type of incoming value. ' +
                                     'value=[{}] vtype=[{}]'.format(type(inc_value), vtype))

        return inc_value

    @staticmethod
    def process_edits(edits, yamlfile):
        '''run through a list of edits and process them one-by-one'''
        results = []
        for edit in edits:
            value = Yedit.parse_value(edit['value'], edit.get('value_type', ''))
            if edit.get('action') == 'update':
                # pylint: disable=line-too-long
                curr_value = Yedit.get_curr_value(
                    Yedit.parse_value(edit.get('curr_value')),
                    edit.get('curr_value_format'))

                rval = yamlfile.update(edit['key'],
                                       value,
                                       edit.get('index'),
                                       curr_value)

            elif edit.get('action') == 'append':
                rval = yamlfile.append(edit['key'], value)

            else:
                rval = yamlfile.put(edit['key'], value)

            if rval[0]:
                results.append({'key': edit['key'], 'edit': rval[1]})

        return {'changed': len(results) > 0, 'results': results}

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params):
        '''perform the idempotent crud operations'''
        yamlfile = Yedit(filename=params['src'],
                         backup=params['backup'],
                         content_type=params['content_type'],
                         backup_ext=params['backup_ext'],
                         separator=params['separator'])

        state = params['state']

        if params['src']:
            rval = yamlfile.load()

            if yamlfile.yaml_dict is None and state != 'present':
                return {'failed': True,
                        'msg': 'Error opening file [{}].  Verify that the '.format(params['src']) +
                               'file exists, that it is has correct permissions, and is valid yaml.'}

        if state == 'list':
            if params['content']:
                content = Yedit.parse_value(params['content'], params['content_type'])
                yamlfile.yaml_dict = content

            if params['key']:
                rval = yamlfile.get(params['key'])

            return {'changed': False, 'result': rval, 'state': state}

        elif state == 'absent':
            if params['content']:
                content = Yedit.parse_value(params['content'], params['content_type'])
                yamlfile.yaml_dict = content

            if params['update']:
                rval = yamlfile.pop(params['key'], params['value'])
            else:
                rval = yamlfile.delete(params['key'], params['index'], params['value'])

            if rval[0] and params['src']:
                yamlfile.write()

            return {'changed': rval[0], 'result': rval[1], 'state': state}

        elif state == 'present':
            # check if content is different than what is in the file
            if params['content']:
                content = Yedit.parse_value(params['content'], params['content_type'])

                # We had no edits to make and the contents are the same
                if yamlfile.yaml_dict == content and \
                   params['value'] is None:
                    return {'changed': False, 'result': yamlfile.yaml_dict, 'state': state}

                yamlfile.yaml_dict = content

            # If we were passed a key, value then
            # we enapsulate it in a list and process it
            # Key, Value passed to the module : Converted to Edits list #
            edits = []
            _edit = {}
            if params['value'] is not None:
                _edit['value'] = params['value']
                _edit['value_type'] = params['value_type']
                _edit['key'] = params['key']

                if params['update']:
                    _edit['action'] = 'update'
                    _edit['curr_value'] = params['curr_value']
                    _edit['curr_value_format'] = params['curr_value_format']
                    _edit['index'] = params['index']

                elif params['append']:
                    _edit['action'] = 'append'

                edits.append(_edit)

            elif params['edits'] is not None:
                edits = params['edits']

            if edits:
                results = Yedit.process_edits(edits, yamlfile)

                # if there were changes and a src provided to us we need to write
                if results['changed'] and params['src']:
                    yamlfile.write()

                return {'changed': results['changed'], 'result': results['results'], 'state': state}

            # no edits to make
            if params['src']:
                # pylint: disable=redefined-variable-type
                rval = yamlfile.write()
                return {'changed': rval[0],
                        'result': rval[1],
                        'state': state}

            # We were passed content but no src, key or value, or edits.  Return contents in memory
            return {'changed': False, 'result': yamlfile.yaml_dict, 'state': state}
        return {'failed': True, 'msg': 'Unkown state passed'}

# -*- -*- -*- End included fragment: ../../lib_utils/src/class/yedit.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: lib/base.py -*- -*- -*-
# pylint: disable=too-many-lines
# noqa: E301,E302,E303,T001


class OpenShiftCLIError(Exception):
    '''Exception class for openshiftcli'''
    pass


ADDITIONAL_PATH_LOOKUPS = ['/usr/local/bin', os.path.expanduser('~/bin')]


def locate_oc_binary():
    ''' Find and return oc binary file '''
    # https://github.com/openshift/openshift-ansible/issues/3410
    # oc can be in /usr/local/bin in some cases, but that may not
    # be in $PATH due to ansible/sudo
    paths = os.environ.get("PATH", os.defpath).split(os.pathsep) + ADDITIONAL_PATH_LOOKUPS

    oc_binary = 'oc'

    # Use shutil.which if it is available, otherwise fallback to a naive path search
    try:
        which_result = shutil.which(oc_binary, path=os.pathsep.join(paths))
        if which_result is not None:
            oc_binary = which_result
    except AttributeError:
        for path in paths:
            if os.path.exists(os.path.join(path, oc_binary)):
                oc_binary = os.path.join(path, oc_binary)
                break

    return oc_binary


# pylint: disable=too-few-public-methods
class OpenShiftCLI(object):
    ''' Class to wrap the command line tools '''
    def __init__(self,
                 namespace,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False,
                 all_namespaces=False):
        ''' Constructor for OpenshiftCLI '''
        self.namespace = namespace
        self.verbose = verbose
        self.kubeconfig = Utils.create_tmpfile_copy(kubeconfig)
        self.all_namespaces = all_namespaces
        self.oc_binary = locate_oc_binary()

    # Pylint allows only 5 arguments to be passed.
    # pylint: disable=too-many-arguments
    def _replace_content(self, resource, rname, content, edits=None, force=False, sep='.'):
        ''' replace the current object with the content '''
        res = self._get(resource, rname)
        if not res['results']:
            return res

        fname = Utils.create_tmpfile(rname + '-')

        yed = Yedit(fname, res['results'][0], separator=sep)
        updated = False

        if content is not None:
            changes = []
            for key, value in content.items():
                changes.append(yed.put(key, value))

            if any([change[0] for change in changes]):
                updated = True

        elif edits is not None:
            results = Yedit.process_edits(edits, yed)

            if results['changed']:
                updated = True

        if updated:
            yed.write()
            atexit.register(Utils.cleanup, [fname])

            return self._replace(fname, force)

        return {'returncode': 0, 'updated': False}

    def _replace(self, fname, force=False):
        '''replace the current object with oc replace'''
        # We are removing the 'resourceVersion' to handle
        # a race condition when modifying oc objects
        yed = Yedit(fname)
        results = yed.delete('metadata.resourceVersion')
        if results[0]:
            yed.write()

        cmd = ['replace', '-f', fname]
        if force:
            cmd.append('--force')
        return self.openshift_cmd(cmd)

    def _create_from_content(self, rname, content):
        '''create a temporary file and then call oc create on it'''
        fname = Utils.create_tmpfile(rname + '-')
        yed = Yedit(fname, content=content)
        yed.write()

        atexit.register(Utils.cleanup, [fname])

        return self._create(fname)

    def _create(self, fname):
        '''call oc create on a filename'''
        return self.openshift_cmd(['create', '-f', fname])

    def _delete(self, resource, name=None, selector=None):
        '''call oc delete on a resource'''
        cmd = ['delete', resource]
        if selector is not None:
            cmd.append('--selector={}'.format(selector))
        elif name is not None:
            cmd.append(name)
        else:
            raise OpenShiftCLIError('Either name or selector is required when calling delete.')

        return self.openshift_cmd(cmd)

    def _process(self, template_name, create=False, params=None, template_data=None):  # noqa: E501
        '''process a template

           template_name: the name of the template to process
           create: whether to send to oc create after processing
           params: the parameters for the template
           template_data: the incoming template's data; instead of a file
        '''
        cmd = ['process']
        if template_data:
            cmd.extend(['-f', '-'])
        else:
            cmd.append(template_name)
        if params:
            param_str = ["{}={}".format(key, str(value).replace("'", r'"')) for key, value in params.items()]
            cmd.append('-p')
            cmd.extend(param_str)

        results = self.openshift_cmd(cmd, output=True, input_data=template_data)

        if results['returncode'] != 0 or not create:
            return results

        fname = Utils.create_tmpfile(template_name + '-')
        yed = Yedit(fname, results['results'])
        yed.write()

        atexit.register(Utils.cleanup, [fname])

        return self.openshift_cmd(['create', '-f', fname])

    def _get(self, resource, name=None, selector=None, field_selector=None):
        '''return a resource by name '''
        cmd = ['get', resource]

        if selector is not None:
            cmd.append('--selector={}'.format(selector))

        if field_selector is not None:
            cmd.append('--field-selector={}'.format(field_selector))

        # Name cannot be used with selector or field_selector.
        if selector is None and field_selector is None and name is not None:
            cmd.append(name)

        cmd.extend(['-o', 'json'])

        rval = self.openshift_cmd(cmd, output=True)

        # Ensure results are retuned in an array
        if 'items' in rval:
            rval['results'] = rval['items']
        elif not isinstance(rval['results'], list):
            rval['results'] = [rval['results']]

        return rval

    def _schedulable(self, node=None, selector=None, schedulable=True):
        ''' perform oadm manage-node scheduable '''
        cmd = ['manage-node']
        if node:
            cmd.extend(node)
        else:
            cmd.append('--selector={}'.format(selector))

        cmd.append('--schedulable={}'.format(schedulable))

        return self.openshift_cmd(cmd, oadm=True, output=True, output_type='raw')  # noqa: E501

    def _list_pods(self, node=None, selector=None, pod_selector=None):
        ''' perform oadm list pods

            node: the node in which to list pods
            selector: the label selector filter if provided
            pod_selector: the pod selector filter if provided
        '''
        cmd = ['manage-node']
        if node:
            cmd.extend(node)
        else:
            cmd.append('--selector={}'.format(selector))

        if pod_selector:
            cmd.append('--pod-selector={}'.format(pod_selector))

        cmd.extend(['--list-pods', '-o', 'json'])

        return self.openshift_cmd(cmd, oadm=True, output=True, output_type='raw')

    # pylint: disable=too-many-arguments
    def _evacuate(self, node=None, selector=None, pod_selector=None, dry_run=False, grace_period=None, force=False):
        ''' perform oadm manage-node evacuate '''
        cmd = ['manage-node']
        if node:
            cmd.extend(node)
        else:
            cmd.append('--selector={}'.format(selector))

        if dry_run:
            cmd.append('--dry-run')

        if pod_selector:
            cmd.append('--pod-selector={}'.format(pod_selector))

        if grace_period:
            cmd.append('--grace-period={}'.format(int(grace_period)))

        if force:
            cmd.append('--force')

        cmd.append('--evacuate')

        return self.openshift_cmd(cmd, oadm=True, output=True, output_type='raw')

    def _version(self):
        ''' return the openshift version'''
        return self.openshift_cmd(['version'], output=True, output_type='raw')

    def _import_image(self, url=None, name=None, tag=None):
        ''' perform image import '''
        cmd = ['import-image']

        image = '{0}'.format(name)
        if tag:
            image += ':{0}'.format(tag)

        cmd.append(image)

        if url:
            cmd.append('--from={0}/{1}'.format(url, image))

        cmd.append('-n{0}'.format(self.namespace))

        cmd.append('--confirm')
        return self.openshift_cmd(cmd)

    def _run(self, cmds, input_data):
        ''' Actually executes the command. This makes mocking easier. '''
        curr_env = os.environ.copy()
        curr_env.update({'KUBECONFIG': self.kubeconfig})
        proc = subprocess.Popen(cmds,
                                stdin=subprocess.PIPE,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE,
                                env=curr_env)

        stdout, stderr = proc.communicate(input_data)

        return proc.returncode, stdout.decode('utf-8'), stderr.decode('utf-8')

    # pylint: disable=too-many-arguments,too-many-branches
    def openshift_cmd(self, cmd, oadm=False, output=False, output_type='json', input_data=None):
        '''Base command for oc '''
        cmds = [self.oc_binary]

        if oadm:
            cmds.append('adm')

        cmds.extend(cmd)

        if self.all_namespaces:
            cmds.extend(['--all-namespaces'])
        elif self.namespace is not None and self.namespace.lower() not in ['none', 'emtpy']:  # E501
            cmds.extend(['-n', self.namespace])

        if self.verbose:
            print(' '.join(cmds))

        try:
            returncode, stdout, stderr = self._run(cmds, input_data)
        except OSError as ex:
            returncode, stdout, stderr = 1, '', 'Failed to execute {}: {}'.format(subprocess.list2cmdline(cmds), ex)

        rval = {"returncode": returncode,
                "cmd": ' '.join(cmds)}

        if output_type == 'json':
            rval['results'] = {}
            if output and stdout:
                try:
                    rval['results'] = json.loads(stdout)
                except ValueError as verr:
                    if "No JSON object could be decoded" in verr.args:
                        rval['err'] = verr.args
        elif output_type == 'raw':
            rval['results'] = stdout if output else ''

        if self.verbose:
            print("STDOUT: {0}".format(stdout))
            print("STDERR: {0}".format(stderr))

        if 'err' in rval or returncode != 0:
            rval.update({"stderr": stderr,
                         "stdout": stdout})

        return rval


class Utils(object):  # pragma: no cover
    ''' utilities for openshiftcli modules '''

    @staticmethod
    def _write(filename, contents):
        ''' Actually write the file contents to disk. This helps with mocking. '''

        with open(filename, 'w') as sfd:
            sfd.write(str(contents))

    @staticmethod
    def create_tmp_file_from_contents(rname, data, ftype='yaml'):
        ''' create a file in tmp with name and contents'''

        tmp = Utils.create_tmpfile(prefix=rname)

        if ftype == 'yaml':
            # AUDIT:no-member makes sense here due to ruamel.YAML/PyYAML usage
            # pylint: disable=no-member
            if hasattr(yaml, 'RoundTripDumper'):
                Utils._write(tmp, yaml.dump(data, Dumper=yaml.RoundTripDumper))
            else:
                Utils._write(tmp, yaml.safe_dump(data, default_flow_style=False))

        elif ftype == 'json':
            Utils._write(tmp, json.dumps(data))
        else:
            Utils._write(tmp, data)

        # Register cleanup when module is done
        atexit.register(Utils.cleanup, [tmp])
        return tmp

    @staticmethod
    def create_tmpfile_copy(inc_file):
        '''create a temporary copy of a file'''
        tmpfile = Utils.create_tmpfile('lib_openshift-')
        Utils._write(tmpfile, open(inc_file).read())

        # Cleanup the tmpfile
        atexit.register(Utils.cleanup, [tmpfile])

        return tmpfile

    @staticmethod
    def create_tmpfile(prefix='tmp'):
        ''' Generates and returns a temporary file name '''

        with tempfile.NamedTemporaryFile(prefix=prefix, delete=False) as tmp:
            return tmp.name

    @staticmethod
    def create_tmp_files_from_contents(content, content_type=None):
        '''Turn an array of dict: filename, content into a files array'''
        if not isinstance(content, list):
            content = [content]
        files = []
        for item in content:
            path = Utils.create_tmp_file_from_contents(item['path'] + '-',
                                                       item['data'],
                                                       ftype=content_type)
            files.append({'name': os.path.basename(item['path']),
                          'path': path})
        return files

    @staticmethod
    def cleanup(files):
        '''Clean up on exit '''
        for sfile in files:
            if os.path.exists(sfile):
                if os.path.isdir(sfile):
                    shutil.rmtree(sfile)
                elif os.path.isfile(sfile):
                    os.remove(sfile)

    @staticmethod
    def exists(results, _name):
        ''' Check to see if the results include the name '''
        if not results:
            return False

        if Utils.find_result(results, _name):
            return True

        return False

    @staticmethod
    def find_result(results, _name):
        ''' Find the specified result by name'''
        rval = None
        for result in results:
            if 'metadata' in result and result['metadata']['name'] == _name:
                rval = result
                break

        return rval

    @staticmethod
    def get_resource_file(sfile, sfile_type='yaml'):
        ''' return the service file '''
        contents = None
        with open(sfile) as sfd:
            contents = sfd.read()

        if sfile_type == 'yaml':
            # AUDIT:no-member makes sense here due to ruamel.YAML/PyYAML usage
            # pylint: disable=no-member
            if hasattr(yaml, 'RoundTripLoader'):
                contents = yaml.load(contents, yaml.RoundTripLoader)
            else:
                contents = yaml.safe_load(contents)
        elif sfile_type == 'json':
            contents = json.loads(contents)

        return contents

    @staticmethod
    def filter_versions(stdout):
        ''' filter the oc version output '''

        version_dict = {}
        version_search = ['oc', 'openshift', 'kubernetes']

        for line in stdout.strip().split('\n'):
            for term in version_search:
                if not line:
                    continue
                if line.startswith(term):
                    version_dict[term] = line.split()[-1]

        # horrible hack to get openshift version in Openshift 3.2
        #  By default "oc version in 3.2 does not return an "openshift" version
        if "openshift" not in version_dict:
            version_dict["openshift"] = version_dict["oc"]

        return version_dict

    @staticmethod
    def add_custom_versions(versions):
        ''' create custom versions strings '''

        versions_dict = {}

        for tech, version in versions.items():
            # clean up "-" from version
            if "-" in version:
                version = version.split("-")[0]

            if version.startswith('v'):
                version = version[1:]  # Remove the 'v' prefix
                versions_dict[tech + '_numeric'] = version.split('+')[0]
                # "3.3.0.33" is what we have, we want "3.3"
                versions_dict[tech + '_short'] = "{}.{}".format(*version.split('.'))

        return versions_dict

    @staticmethod
    def openshift_installed():
        ''' check if openshift is installed '''
        import rpm

        transaction_set = rpm.TransactionSet()
        rpmquery = transaction_set.dbMatch("name", "atomic-openshift")

        return rpmquery.count() > 0

    # Disabling too-many-branches.  This is a yaml dictionary comparison function
    # pylint: disable=too-many-branches,too-many-return-statements,too-many-statements
    @staticmethod
    def check_def_equal(user_def, result_def, skip_keys=None, debug=False):
        ''' Given a user defined definition, compare it with the results given back by our query.  '''

        # Currently these values are autogenerated and we do not need to check them
        skip = ['metadata', 'status']
        if skip_keys:
            skip.extend(skip_keys)

        for key, value in result_def.items():
            if key in skip:
                continue

            # Both are lists
            if isinstance(value, list):
                if key not in user_def:
                    if debug:
                        print('User data does not have key [%s]' % key)
                        print('User data: %s' % user_def)
                    return False

                if not isinstance(user_def[key], list):
                    if debug:
                        print('user_def[key] is not a list key=[%s] user_def[key]=%s' % (key, user_def[key]))
                    return False

                if len(user_def[key]) != len(value):
                    if debug:
                        print("List lengths are not equal.")
                        print("key=[%s]: user_def[%s] != value[%s]" % (key, len(user_def[key]), len(value)))
                        print("user_def: %s" % user_def[key])
                        print("value: %s" % value)
                    return False

                for values in zip(user_def[key], value):
                    if isinstance(values[0], dict) and isinstance(values[1], dict):
                        if debug:
                            print('sending list - list')
                            print(type(values[0]))
                            print(type(values[1]))
                        result = Utils.check_def_equal(values[0], values[1], skip_keys=skip_keys, debug=debug)
                        if not result:
                            print('list compare returned false')
                            return False

                    elif value != user_def[key]:
                        if debug:
                            print('value should be identical')
                            print(user_def[key])
                            print(value)
                        return False

            # recurse on a dictionary
            elif isinstance(value, dict):
                if key not in user_def:
                    if debug:
                        print("user_def does not have key [%s]" % key)
                    return False
                if not isinstance(user_def[key], dict):
                    if debug:
                        print("dict returned false: not instance of dict")
                    return False

                # before passing ensure keys match
                api_values = set(value.keys()) - set(skip)
                user_values = set(user_def[key].keys()) - set(skip)
                if api_values != user_values:
                    if debug:
                        print("keys are not equal in dict")
                        print(user_values)
                        print(api_values)
                    return False

                result = Utils.check_def_equal(user_def[key], value, skip_keys=skip_keys, debug=debug)
                if not result:
                    if debug:
                        print("dict returned false")
                        print(result)
                    return False

            # Verify each key, value pair is the same
            else:
                if key not in user_def or value != user_def[key]:
                    if debug:
                        print("value not equal; user_def does not have key")
                        print(key)
                        print(value)
                        if key in user_def:
                            print(user_def[key])
                    return False

        if debug:
            print('returning true')
        return True

class OpenShiftCLIConfig(object):
    '''Generic Config'''
    def __init__(self, rname, namespace, kubeconfig, options):
        self.kubeconfig = kubeconfig
        self.name = rname
        self.namespace = namespace
        self._options = options

    @property
    def config_options(self):
        ''' return config options '''
        return self._options

    def to_option_list(self, ascommalist=''):
        '''return all options as a string
           if ascommalist is set to the name of a key, and
           the value of that key is a dict, format the dict
           as a list of comma delimited key=value pairs'''
        return self.stringify(ascommalist)

    def stringify(self, ascommalist=''):
        ''' return the options hash as cli params in a string
            if ascommalist is set to the name of a key, and
            the value of that key is a dict, format the dict
            as a list of comma delimited key=value pairs '''
        rval = []
        for key in sorted(self.config_options.keys()):
            data = self.config_options[key]
            if data['include'] \
               and (data['value'] is not None or isinstance(data['value'], int)):
                if key == ascommalist:
                    val = ','.join(['{}={}'.format(kk, vv) for kk, vv in sorted(data['value'].items())])
                else:
                    val = data['value']
                rval.append('--{}={}'.format(key.replace('_', '-'), val))

        return rval


# -*- -*- -*- End included fragment: lib/base.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: lib/group.py -*- -*- -*-


class GroupConfig(object):
    ''' Handle route options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 kubeconfig):
        ''' constructor for handling group options '''
        self.kubeconfig = kubeconfig
        self.name = sname
        self.namespace = namespace
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' return a service as a dict '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'Group'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        self.data['users'] = None


# pylint: disable=too-many-instance-attributes
class Group(Yedit):
    ''' Class to wrap the oc command line tools '''
    kind = 'group'

    def __init__(self, content):
        '''Group constructor'''
        super(Group, self).__init__(content=content)

# -*- -*- -*- End included fragment: lib/group.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: class/oc_group.py -*- -*- -*-


class OCGroup(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'group'

    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCGroup '''
        super(OCGroup, self).__init__(config.namespace, config.kubeconfig)
        self.config = config
        self.namespace = config.namespace
        self._group = None

    @property
    def group(self):
        ''' property function service'''
        if not self._group:
            self.get()
        return self._group

    @group.setter
    def group(self, data):
        ''' setter function for yedit var '''
        self._group = data

    def exists(self):
        ''' return whether a group exists '''
        if self.group:
            return True

        return False

    def get(self):
        '''return group information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.group = Group(content=result['results'][0])
        elif 'groups.user.openshift.io \"{}\" not found'.format(self.config.name) in result['stderr']:
            result['returncode'] = 0
            result['results'] = [{}]

        return result

    def delete(self):
        '''delete the object'''
        return self._delete(self.kind, self.config.name)

    def create(self):
        '''create the object'''
        return self._create_from_content(self.config.name, self.config.data)

    def update(self):
        '''update the object'''
        return self._replace_content(self.kind, self.config.name, self.config.data)

    def needs_update(self):
        ''' verify an update is needed '''
        return not Utils.check_def_equal(self.config.data, self.group.yaml_dict, skip_keys=['users'], debug=True)

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode=False):
        '''run the oc_group module'''

        gconfig = GroupConfig(params['name'],
                              params['namespace'],
                              params['kubeconfig'],
                             )
        oc_group = OCGroup(gconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_group.get()

        if api_rval['returncode'] != 0:
            return {'failed': True, 'msg': api_rval}

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_group.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

                api_rval = oc_group.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_group.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                # Create it here
                api_rval = oc_group.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_group.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_group.needs_update():
                api_rval = oc_group.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_group.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'Unknown state passed. {}'.format(state)}

# -*- -*- -*- End included fragment: class/oc_group.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: ansible/oc_group.py -*- -*- -*-

#pylint: disable=too-many-branches
def main():
    '''
    ansible oc module for group
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, type='str'),
            namespace=dict(default='default', type='str'),
            # addind users to a group is handled through the oc_users module
            #users=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
    )

    rval = OCGroup.run_ansible(module.params, module.check_mode)

    if 'failed' in rval:
        return module.fail_json(**rval)

    return module.exit_json(**rval)

if __name__ == '__main__':
    main()

# -*- -*- -*- End included fragment: ansible/oc_group.py -*- -*- -*-
