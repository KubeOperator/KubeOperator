'''
 Unit tests for the FakeOpenSSL classes
'''
import os
import subprocess
import sys

import pytest

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'library'))
sys.path.insert(1, MODULE_PATH)

# pylint: disable=import-error,wrong-import-position,missing-docstring
# pylint: disable=invalid-name,redefined-outer-name
from openshift_cert_expiry import FakeOpenSSLCertificate  # noqa: E402


@pytest.fixture(scope='module')
def fake_valid_cert(valid_cert):
    cmd = ['openssl', 'x509', '-in', str(valid_cert['cert_file']), '-text',
           '-nameopt', 'oneline']
    cert = subprocess.check_output(cmd)
    return FakeOpenSSLCertificate(cert.decode('utf8'))


def test_not_after(valid_cert, fake_valid_cert):
    ''' Validate value returned back from get_notAfter() '''
    real_cert = valid_cert['cert']

    # Internal representation of pyOpenSSL is bytes, while FakeOpenSSLCertificate
    # is text, so decode the result from pyOpenSSL prior to comparing
    assert real_cert.get_notAfter().decode('utf8') == fake_valid_cert.get_notAfter()


def test_serial(valid_cert, fake_valid_cert):
    ''' Validate value returned back form get_serialnumber() '''
    real_cert = valid_cert['cert']
    assert real_cert.get_serial_number() == fake_valid_cert.get_serial_number()


def test_get_subject(valid_cert, fake_valid_cert):
    ''' Validate the certificate subject '''

    # Gather the subject components and create a list of colon separated strings.
    # Since the internal representation of pyOpenSSL uses bytes, we need to decode
    # the results before comparing.
    c_subjects = valid_cert['cert'].get_subject().get_components()
    c_subj = ', '.join(['{}:{}'.format(x.decode('utf8'), y.decode('utf8')) for x, y in c_subjects])
    f_subjects = fake_valid_cert.get_subject().get_components()
    f_subj = ', '.join(['{}:{}'.format(x, y) for x, y in f_subjects])
    assert c_subj == f_subj


def get_san_extension(cert):
    # Internal representation of pyOpenSSL is bytes, while FakeOpenSSLCertificate
    # is text, so we need to set the value to search for accordingly.
    if isinstance(cert, FakeOpenSSLCertificate):
        san_short_name = 'subjectAltName'
    else:
        san_short_name = b'subjectAltName'

    for i in range(cert.get_extension_count()):
        ext = cert.get_extension(i)
        if ext.get_short_name() == san_short_name:
            # return the string representation to compare the actual SAN
            # values instead of the data types
            return str(ext)

    return None


def test_subject_alt_names(valid_cert, fake_valid_cert):
    real_cert = valid_cert['cert']

    san = get_san_extension(real_cert)
    f_san = get_san_extension(fake_valid_cert)

    assert san == f_san

    # If there are either dns or ip sans defined, verify common_name present
    if valid_cert['ip'] or valid_cert['dns']:
        assert 'DNS:' + valid_cert['common_name'] in f_san

    # Verify all ip sans are present
    for ip in valid_cert['ip']:
        assert 'IP Address:' + ip in f_san

    # Verify all dns sans are present
    for name in valid_cert['dns']:
        assert 'DNS:' + name in f_san
