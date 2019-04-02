#!/usr/bin/env python
'''
  Generate the openshift-ansible/roles/lib_openshift_cli/library/ modules.
'''

import argparse
import os
import re
import yaml
import six

OPENSHIFT_ANSIBLE_PATH = os.path.dirname(os.path.realpath(__file__))
OPENSHIFT_ANSIBLE_SOURCES_PATH = os.path.join(OPENSHIFT_ANSIBLE_PATH, 'sources.yml')  # noqa: E501
LIBRARY = os.path.join(OPENSHIFT_ANSIBLE_PATH, '..', 'library/')
SKIP_COVERAGE_PATTERN = [re.compile('class Yedit.*$'),
                         re.compile('class Utils.*$')]
PRAGMA_STRING = '  # pragma: no cover'


class GenerateAnsibleException(Exception):
    '''General Exception for generate function'''
    pass


def parse_args():
    '''parse arguments to generate'''
    parser = argparse.ArgumentParser(description="Generate ansible modules.")
    parser.add_argument('--verify', action='store_true', default=False,
                        help='Verify library code matches the generated code.')

    return parser.parse_args()


def fragment_banner(fragment_path, side, data):
    """Generate a banner to wrap around file fragments

:param string fragment_path: A path to a module fragment
:param string side: ONE OF: "header", "footer"
:param StringIO data: A StringIO object to write the banner to
"""
    side_msg = {
        "header": "Begin included fragment: {}",
        "footer": "End included fragment: {}"
    }
    annotation = side_msg[side].format(fragment_path)

    banner = """
# -*- -*- -*- {} -*- -*- -*-
""".format(annotation)

    # Why skip?
    #
    # * 'generated' - This is the head of the script, we don't want to
    #   put comments before the #!shebang
    #
    # * 'license' - Wrapping this just seemed like gratuitous extra
    if ("generated" not in fragment_path) and ("license" not in fragment_path):
        data.write(banner)

    # Make it self-contained testable
    return banner


def generate(parts):
    '''generate the source code for the ansible modules

:param Array parts: An array of paths (strings) to module fragments
    '''

    data = six.StringIO()
    for fpart in parts:
        # first line is pylint disable so skip it
        with open(os.path.join(OPENSHIFT_ANSIBLE_PATH, fpart)) as pfd:
            fragment_banner(fpart, "header", data)
            for idx, line in enumerate(pfd):
                if idx in [0, 1] and 'flake8: noqa' in line or 'pylint: skip-file' in line:  # noqa: E501
                    continue

                for skip in SKIP_COVERAGE_PATTERN:
                    if re.match(skip, line):
                        line = line.strip()
                        line += PRAGMA_STRING + os.linesep

                data.write(line)

            fragment_banner(fpart, "footer", data)
    return data


def get_sources():
    '''return the path to the generate sources'''
    return yaml.load(open(OPENSHIFT_ANSIBLE_SOURCES_PATH).read())


def verify():
    '''verify if the generated code matches the library code'''
    for fname, parts in get_sources().items():
        data = generate(parts)
        fname = os.path.join(LIBRARY, fname)
        if not open(fname).read() == data.getvalue():
            raise GenerateAnsibleException('Generated content does not match for %s' % fname)


def main():
    ''' combine the necessary files to create the ansible module '''
    args = parse_args()
    if args.verify:
        verify()

    for fname, parts in get_sources().items():
        data = generate(parts)
        fname = os.path.join(LIBRARY, fname)
        with open(fname, 'w') as afd:
            afd.seek(0)
            afd.write(data.getvalue())


if __name__ == '__main__':
    main()
