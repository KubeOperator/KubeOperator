'''
 Unit tests for the load_and_handle_cert method
'''
import datetime
import os
import sys

import pytest

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'library'))
sys.path.insert(1, MODULE_PATH)

# pylint: disable=import-error,wrong-import-position,missing-docstring
# pylint: disable=invalid-name,redefined-outer-name
import openshift_cert_expiry  # noqa: E402

# TODO: More testing on the results of the load_and_handle_cert function
# could be implemented here as well, such as verifying subjects
# match up.


@pytest.fixture(params=['OpenSSLCertificate', 'FakeOpenSSLCertificate'])
def loaded_cert(request, valid_cert):
    """ parameterized fixture to provide load_and_handle_cert results
        for both OpenSSL and FakeOpenSSL parsed certificates
    """
    now = datetime.datetime.now()

    openshift_cert_expiry.HAS_OPENSSL = request.param == 'OpenSSLCertificate'

    # valid_cert['cert_file'] is a `py.path.LocalPath` object and
    # provides a read_text() method for reading the file contents.
    cert_string = valid_cert['cert_file'].read_text('utf8')

    (subject,
     expiry_date,
     time_remaining,
     serial) = openshift_cert_expiry.load_and_handle_cert(cert_string, now)

    return {
        'now': now,
        'subject': subject,
        'expiry_date': expiry_date,
        'time_remaining': time_remaining,
        'serial': serial,
    }


def test_serial(loaded_cert, valid_cert):
    """Params:

    * `loaded_cert` comes from the `loaded_cert` fixture in this file
    * `valid_cert` comes from the 'valid_cert' fixture in conftest.py
    """
    valid_cert_serial = valid_cert['cert'].get_serial_number()
    assert loaded_cert['serial'] == valid_cert_serial


def test_expiry(loaded_cert):
    """Params:

    * `loaded_cert` comes from the `loaded_cert` fixture in this file
    """
    expiry_date = loaded_cert['expiry_date']
    time_remaining = loaded_cert['time_remaining']
    now = loaded_cert['now']
    assert expiry_date == now + time_remaining
