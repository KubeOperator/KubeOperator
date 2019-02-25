'''
 Unit tests for repoquery
'''

import os
import sys
import unittest
import mock

# Removing invalid variable names for tests so that I can
# keep them brief
# pylint: disable=invalid-name,no-name-in-module
# Disable import-error b/c our libraries aren't loaded in jenkins
# pylint: disable=import-error,wrong-import-position
# place class in our python path
module_path = os.path.join('/'.join(os.path.realpath(__file__).split('/')[:-4]), 'library')  # noqa: E501
sys.path.insert(0, module_path)
from repoquery import Repoquery  # noqa: E402


class RepoQueryTest(unittest.TestCase):
    '''
     Test class for RepoQuery
    '''

    @mock.patch('repoquery._run')
    def test_querying_a_package(self, mock_cmd):
        ''' Testing querying a package '''

        # Arrange

        # run_ansible input parameters
        params = {
            'state': 'list',
            'name': 'bash',
            'query_type': 'repos',
            'verbose': False,
            'show_duplicates': False,
            'match_version': None,
            'ignore_excluders': False,
        }

        valid_stderr = '''Repo rhel-7-server-extras-rpms forced skip_if_unavailable=True due to: /etc/pki/entitlement/3268107132875399464-key.pem
        Repo rhel-7-server-rpms forced skip_if_unavailable=True due to: /etc/pki/entitlement/4128505182875899164-key.pem'''  # not real

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            (0, b'4.2.46|21.el7_3|x86_64|rhel-7-server-rpms|4.2.46-21.el7_3', valid_stderr),  # first call to the mock
        ]

        # Act
        results = Repoquery.run_ansible(params, False)

        # Assert
        self.assertEqual(results['state'], 'list')
        self.assertFalse(results['changed'])
        self.assertTrue(results['results']['package_found'])
        self.assertEqual(results['results']['returncode'], 0)
        self.assertEqual(results['results']['package_name'], 'bash')
        self.assertEqual(results['results']['versions'], {'latest_full': '4.2.46-21.el7_3',
                                                          'available_versions': ['4.2.46'],
                                                          'available_versions_full': ['4.2.46-21.el7_3'],
                                                          'latest': '4.2.46'})

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['/usr/bin/repoquery', '--plugins', '--quiet', '--pkgnarrow=repos', '--queryformat=%{version}|%{release}|%{arch}|%{repo}|%{version}-%{release}', 'bash']),
        ])
