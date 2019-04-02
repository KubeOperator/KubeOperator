'''
 Unit tests for the master_check_paths_in_config action plugin
'''
import os
import sys

from ansible import errors
import pytest


MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'action_plugins'))
sys.path.insert(1, MODULE_PATH)

# pylint: disable=import-error,wrong-import-position,missing-docstring
# pylint: disable=invalid-name,redefined-outer-name
import master_check_paths_in_config  # noqa: E402


@pytest.fixture()
def loaded_config():
    """return testing master config"""
    data = {
        'apiVersion': 'v1',
        'oauthConfig':
        {'identityProviders':
            ['1', '2', '/this/will/fail']},
        'fake_top_item':
        {'fake_item':
            {'fake_item2':
                ["some string",
                    {"fake_item3":
                        ["some string 2",
                            {"fake_item4":
                                {"some_key": "deeply_nested_string"}}]}]}}
    }
    return data


def test_pop_migrated(loaded_config):
    """Params:

    * `loaded_config` comes from the `loaded_config` fixture in this file
    """
    # Ensure we actually loaded a valid config
    assert loaded_config['apiVersion'] == 'v1'

    # Test that migrated key is removed
    master_check_paths_in_config.pop_migrated_fields(loaded_config)
    assert loaded_config['oauthConfig'] is not None
    assert loaded_config['oauthConfig'].get('identityProviders') is None


def test_walk_mapping(loaded_config):
    """Params:
    * `loaded_config` comes from the `loaded_config` fixture in this file
    """
    # Ensure we actually loaded a valid config
    fake_top_item = loaded_config['fake_top_item']
    stc = []
    expected_keys = ("some string", "some string 2", "deeply_nested_string")

    # Test that we actually extract all the strings from complicated nested
    # structures
    master_check_paths_in_config.walk_mapping(fake_top_item, stc)
    assert len(stc) == 3
    for item in expected_keys:
        assert item in stc


def test_check_strings():
    stc_good = ('/etc/origin/master/good', 'some/child/dir')
    # This should not raise
    master_check_paths_in_config.check_strings(stc_good)

    # This is a string we should alert on
    stc_bad = ('goodfile.txt', '/root/somefile')
    with pytest.raises(errors.AnsibleModuleError):
        master_check_paths_in_config.check_strings(stc_bad)

    stc_bad_relative = ('goodfile.txt', '../node/otherfile')
    with pytest.raises(errors.AnsibleModuleError):
        master_check_paths_in_config.check_strings(stc_bad_relative)
