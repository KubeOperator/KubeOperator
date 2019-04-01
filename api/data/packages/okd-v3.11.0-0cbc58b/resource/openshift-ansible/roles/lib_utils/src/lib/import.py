# flake8: noqa
# pylint: skip-file

# pylint: disable=wrong-import-order,wrong-import-position,unused-import

from __future__ import print_function  # noqa: F401
import copy  # noqa: F401
import fcntl  # noqa: F401
import json   # noqa: F401
import os  # noqa: F401
import re  # noqa: F401
import shutil  # noqa: F401
import tempfile  # noqa: F401
import time  # noqa: F401

try:
    import ruamel.yaml as yaml  # noqa: F401
except ImportError:
    import yaml  # noqa: F401

from ansible.module_utils.basic import AnsibleModule
