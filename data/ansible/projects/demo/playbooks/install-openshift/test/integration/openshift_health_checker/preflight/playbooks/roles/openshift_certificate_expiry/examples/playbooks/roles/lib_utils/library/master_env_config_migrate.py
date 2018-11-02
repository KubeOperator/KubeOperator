#!/usr/bin/env python
# pylint: disable=missing-docstring
#
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

try:
    # configparser is available in python 2.7 backports, but that package
    # might not be installed.
    import ConfigParser as configparser
except ImportError:
    import configparser
import re
import sys
import os

from ansible.module_utils.basic import AnsibleModule


DOCUMENTATION = '''
---
module: master_env_config_migrate

short_description: Migrates an environment file from one location to another.

version_added: "2.4"

description:
    - Ensures that an environment file is properly migrated and values are properly
      quoted.

options:
    src:
        description:
            - This is the original file on remote host.
        required: true
    dest:
        description:
            - This is the output location.
        required: true

author:
    - "Michael Gugino <mgugino@redhat.com>"
'''


class SectionlessParser(configparser.RawConfigParser):
    # pylint: disable=invalid-name,too-many-locals,too-many-branches,too-many-statements
    # pylint: disable=anomalous-backslash-in-string,raising-bad-type
    """RawConfigParser that allows no sections"""
    # This class originally retrieved from:
    # https://github.com/python/cpython/blob/master/Lib/configparser.py
    # Copyright 2001-2018 Python Software Foundation; All Rights Reserved
    # Modified to allow no sections.
    def _set_proxies(self, sectname):
        """proxies not present in old version"""
        pass

    def optionxform(self, optionstr):
        """Override this method, don't set .lower()"""
        return optionstr

    def _write_section(self, fp, _, section_items, delimiter):
        """Override for formatting"""
        for key, value in section_items:
            if " " in value and "\ " not in value and not value.startswith('"'):
                value = u'"{}"'.format(value)
            if value is not None or not self._allow_no_value:
                value = delimiter + str(value).replace('\n', '\n\t')
            else:
                value = u""
            fp.write(u"{}{}\n".format(key, value))
            fp.write(u"\n")

    # pylint: disable=arguments-differ
    def write(self, fp, space_around_delimiters=True):
        """Ovrride write method"""
        delimiters = ('=', ':')
        if space_around_delimiters:
            d = " {} ".format(delimiters[0])
        else:
            d = delimiters[0]
        for section in self._sections:
            self._write_section(fp, section,
                                self._sections[section].items(), d)

    def _join_multiline_values(self):
        all_sections = self._sections.items()
        for _, options in all_sections:
            for name, val in options.items():
                if isinstance(val, list):
                    val = '\n'.join(val).rstrip()
                options[name] = val

    # pylint: disable=attribute-defined-outside-init
    def _read(self, fp, fpname):
        """Parse a sectionless configuration file."""
        elements_added = set()
        cursect = {}
        sectname = '__none_sect'
        self._sections[sectname] = cursect
        self._set_proxies(sectname)
        self._inline_comment_prefixes = ('#',)
        self._comment_prefixes = ('#',)
        self._empty_lines_in_values = True
        self.NONSPACECRE = re.compile(r"\S")
        self.default_section = 'DEFAULT'
        optname = None
        lineno = 0
        indent_level = 0
        e = None                              # None, or an exception
        for lineno, line in enumerate(fp, start=1):
            comment_start = sys.maxsize
            # strip inline comments
            inline_prefixes = {p: -1 for p in self._inline_comment_prefixes}
            while comment_start == sys.maxsize and inline_prefixes:
                next_prefixes = {}
                for prefix, index in inline_prefixes.items():
                    index = line.find(prefix, index + 1)
                    if index == -1:
                        continue
                    next_prefixes[prefix] = index
                    if index == 0 or (index > 0 and line[index - 1].isspace()):
                        comment_start = min(comment_start, index)
                inline_prefixes = next_prefixes
            # strip full line comments
            for prefix in self._comment_prefixes:
                if line.strip().startswith(prefix):
                    comment_start = 0
                    break
            if comment_start == sys.maxsize:
                comment_start = None
            value = line[:comment_start].strip()
            if not value:
                if self._empty_lines_in_values:
                    # add empty line to the value, but only if there was no
                    # comment on the line
                    if (comment_start is None and
                            cursect is not None and
                            optname and
                            cursect[optname] is not None):
                        cursect[optname].append('')  # newlines added at join
                else:
                    # empty line marks end of value
                    indent_level = sys.maxsize
                continue
            # continuation line?
            first_nonspace = self.NONSPACECRE.search(line)
            cur_indent_level = first_nonspace.start() if first_nonspace else 0
            if (cursect is not None and optname and
                    cur_indent_level > indent_level):
                cursect[optname].append(value)

            # a section header or option header?
            else:
                indent_level = cur_indent_level
                # is it a section header?
                mo = self.SECTCRE.match(value)
                if mo:
                    optname = None
                else:
                    mo = self._optcre.match(value)
                    if mo:
                        optname, _, optval = mo.group('option', 'vi', 'value')
                        if not optname:
                            e = self._handle_error(e, fpname, lineno, line)
                        optname = self.optionxform(optname.rstrip())
                        elements_added.add((sectname, optname))
                        if optval is not None:
                            optval = optval.strip()
                            cursect[optname] = [optval]
                        else:
                            # valueless option handling
                            cursect[optname] = None
                    else:
                        e = self._handle_error(e, fpname, lineno, line)
        self._join_multiline_values()
        # if any parsing errors occurred, raise an exception
        if e:
            raise e


def create_file(src, dest):
    '''Create the dest file from src file'''
    config = SectionlessParser()
    config.readfp(open(src))
    with open(dest, 'w') as output:
        config.write(output, False)


def run_module():
    '''Run this module'''
    module_args = dict(
        src=dict(required=True, type='path'),
        dest=dict(required=True, type='path'),
    )

    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=False
    )

    # First, create our dest dir if necessary
    dest = module.params['dest']
    src = module.params['src']

    if os.path.exists(dest):
        # Do nothing, output file already in place.
        result = {'changed': False}
        module.exit_json(**result)

    create_file(src, dest)

    result = {'changed': True}
    module.exit_json(**result)


def main():
    run_module()


if __name__ == '__main__':
    main()
