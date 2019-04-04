import os
import sys

import pytest

from ansible.module_utils.basic import AnsibleModule

try:
    # python3, mock is built in.
    from unittest.mock import patch
except ImportError:
    # In python2, mock is installed via pip.
    from mock import patch

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'library'))
sys.path.insert(1, MODULE_PATH)

import oc_csr_approve  # noqa
from oc_csr_approve import CSRapprove # noqa

# base path for text files with sample outputs.
ASSET_PATH = os.path.realpath(os.path.join(__file__, os.pardir, 'test_data'))

RUN_CMD_MOCK = 'ansible.module_utils.basic.AnsibleModule.run_command'


class DummyModule(AnsibleModule):
    def _load_params(self):
        self.params = {}

    def exit_json(*args, **kwargs):
        return 0

    def fail_json(*args, **kwargs):
        raise Exception(kwargs['msg'])


def test_parse_subject_cn():
    subject = 'subject=/C=US/CN=fedora1.openshift.io/L=Raleigh/O=Red Hat/ST=North Carolina/OU=OpenShift\n'
    assert oc_csr_approve.parse_subject_cn(subject) == 'fedora1.openshift.io'

    subject = 'subject=C = US, CN = test.io, L = City, O = Company, ST = State, OU = Dept\n'
    assert oc_csr_approve.parse_subject_cn(subject) == 'test.io'


def test_get_nodes():
    output_file = os.path.join(ASSET_PATH, 'oc_get_nodes.json')
    with open(output_file) as stdoutfile:
        oc_get_nodes_stdout = stdoutfile.read()

    module = DummyModule({})
    approver = CSRapprove(module, 'oc', '/dev/null', [])

    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, oc_get_nodes_stdout, '')
        all_nodes = approver.get_nodes()
    assert all_nodes == ['fedora1.openshift.io', 'fedora2.openshift.io', 'fedora3.openshift.io']


def test_get_csrs():
    module = DummyModule({})
    approver = CSRapprove(module, 'oc', '/dev/null', [])
    output_file = os.path.join(ASSET_PATH, 'oc_csr_approve_pending.json')
    with open(output_file) as stdoutfile:
        oc_get_csr_out = stdoutfile.read()

    # mock oc get csr call to cluster
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, oc_get_csr_out, '')
        csrs = approver.get_csrs()

    assert csrs[0]['kind'] == "CertificateSigningRequest"

    output_file = os.path.join(ASSET_PATH, 'openssl1.txt')
    with open(output_file) as stdoutfile:
        openssl_out = stdoutfile.read()

    # mock openssl req call.
    node_list = ['fedora2.mguginolocal.com']
    approver = CSRapprove(module, 'oc', '/dev/null', node_list)
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, openssl_out, '')
        csr_dict = approver.process_csrs(csrs, "client")
    # actually run openssl req call.
    csr_dict = approver.process_csrs(csrs, "client")
    assert csr_dict['node-csr-TkefytQp8Dz4Xp7uzcw605MocvI0gWuEOGNrHhOjGNQ'] == 'fedora2.mguginolocal.com'


def test_confirm_needed_requests_present():
    module = DummyModule({})
    csr_dict = {'some-csr': 'fedora1.openshift.io'}
    not_found_nodes = ['host1']
    approver = CSRapprove(module, 'oc', '/dev/null', [])
    with pytest.raises(Exception) as err:
        approver.confirm_needed_requests_present(not_found_nodes, csr_dict)
        assert 'Exception: Could not find csr for nodes: host1' in str(err)

    not_found_nodes = ['fedora1.openshift.io']
    # this should complete silently
    approver.confirm_needed_requests_present(not_found_nodes, csr_dict)


def test_approve_csrs():
    module = DummyModule({})
    csr_dict = {'csr-1': 'example.openshift.io'}
    approver = CSRapprove(module, 'oc', '/dev/null', [])
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, 'csr-1 ok', '')
        approver.approve_csrs(csr_dict, 'client')
    assert approver.result['client_approve_results'] == ['csr-1 ok']


def test_get_ready_nodes_server():
    module = DummyModule({})
    nodes_list = ['fedora1.openshift.io']
    approver = CSRapprove(module, 'oc', '/dev/null', nodes_list)
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, 'ok', '')
        ready_nodes_server = approver.get_ready_nodes_server(nodes_list)
    assert ready_nodes_server == ['fedora1.openshift.io']


def test_get_csrs_server():
    module = DummyModule({})
    output_file = os.path.join(ASSET_PATH, 'oc_csr_server_multiple_pends_one_host.json')
    with open(output_file) as stdoutfile:
        oc_get_csr_out = stdoutfile.read()

    approver = CSRapprove(module, 'oc', '/dev/null', [])
    # mock oc get csr call to cluster
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, oc_get_csr_out, '')
        csrs = approver.get_csrs()

    assert csrs[0]['kind'] == "CertificateSigningRequest"

    output_file = os.path.join(ASSET_PATH, 'openssl1.txt')
    with open(output_file) as stdoutfile:
        openssl_out = stdoutfile.read()

    node_list = ['fedora1.openshift.io']
    approver = CSRapprove(module, 'oc', '/dev/null', node_list)
    # mock openssl req call.
    with patch(RUN_CMD_MOCK) as call_mock:
        call_mock.return_value = (0, openssl_out, '')
        csr_dict = approver.process_csrs(csrs, "server")

    # actually run openssl req call.
    node_list = ['fedora2.mguginolocal.com']
    approver = CSRapprove(module, 'oc', '/dev/null', node_list)
    csr_dict = approver.process_csrs(csrs, "server")
    assert csr_dict['csr-2cxkp'] == 'fedora2.mguginolocal.com'


def test_verify_server_csrs():
    module = DummyModule({})
    ready_nodes_server = ['fedora1.openshift.io']
    node_list = ['fedora1.openshift.io']
    approver = CSRapprove(module, 'oc', '/dev/null', node_list)
    with patch('oc_csr_approve.CSRapprove.get_ready_nodes_server') as call_mock:
        call_mock.return_value = ready_nodes_server
        # This should silently return
        approver.verify_server_csrs()

    node_list = ['fedora1.openshift.io', 'fedora2.openshift.io']
    approver = CSRapprove(module, 'oc', '/dev/null', node_list)
    with patch('oc_csr_approve.CSRapprove.get_ready_nodes_server') as call_mock:
        call_mock.return_value = ready_nodes_server
        with pytest.raises(Exception) as err:
            approver.verify_server_csrs()
        assert 'after approving server certs: fedora2.openshift.io' in str(err)


if __name__ == '__main__':
    test_parse_subject_cn()
    test_get_nodes()
    test_get_csrs()
    test_confirm_needed_requests_present()
    test_approve_csrs()
    test_get_ready_nodes_server()
    test_get_csrs_server()
    test_verify_server_csrs()
