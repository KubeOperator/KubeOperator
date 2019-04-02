'''
 Unit tests for oo_filters
'''
import os
import sys

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'filter_plugins'))
sys.path.insert(0, MODULE_PATH)

# pylint: disable=import-error,wrong-import-position,missing-docstring
import oo_filters   # noqa: E402


def test_lib_utils_oo_oreg_image():
    default_url = "quay.io/coreos/etcd:v99"

    oreg_url = "None"
    output_image = oo_filters.lib_utils_oo_oreg_image(default_url, oreg_url)
    assert output_image == default_url

    oreg_url = "example.com/openshift/origin-${component}:${version}"
    expected_output = "example.com/coreos/etcd:v99"
    output_image = oo_filters.lib_utils_oo_oreg_image(default_url, oreg_url)
    assert output_image == expected_output

    oreg_url = "example.com/subdir/openshift/origin-${component}:${version}"
    expected_output = "example.com/subdir/coreos/etcd:v99"
    output_image = oo_filters.lib_utils_oo_oreg_image(default_url, oreg_url)
    assert output_image == expected_output


def main():
    test_lib_utils_oo_oreg_image()


if __name__ == '__main__':
    main()
