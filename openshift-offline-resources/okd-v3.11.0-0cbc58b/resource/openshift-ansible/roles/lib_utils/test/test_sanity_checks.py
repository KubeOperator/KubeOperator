'''
 Unit tests for wildcard
'''
import os
import sys

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'action_plugins'))
sys.path.insert(0, MODULE_PATH)

# pylint: disable=import-error,wrong-import-position,missing-docstring
from sanity_checks import is_registry_match   # noqa: E402


def test_is_registry_match():
    '''
     Test for is_registry_match
    '''
    pat_allowall = "*"
    pat_docker = "docker.io"
    pat_subdomain = "*.example.com"
    pat_matchport = "registry:80"

    assert is_registry_match("docker.io/repo/my", pat_allowall)
    assert is_registry_match("example.com:4000/repo/my", pat_allowall)
    assert is_registry_match("172.192.222.10:4000/a/b/c", pat_allowall)
    assert is_registry_match("https://registry.com", pat_allowall)
    assert is_registry_match("example.com/openshift3/ose-${component}:${version}", pat_allowall)

    assert is_registry_match("docker.io/repo/my", pat_docker)
    assert is_registry_match("docker.io:443/repo/my", pat_docker)
    assert is_registry_match("docker.io/openshift3/ose-${component}:${version}", pat_allowall)
    assert not is_registry_match("example.com:4000/repo/my", pat_docker)
    assert not is_registry_match("index.docker.io/a/b/c", pat_docker)
    assert not is_registry_match("https://registry.com", pat_docker)
    assert not is_registry_match("example.com/openshift3/ose-${component}:${version}", pat_docker)

    assert is_registry_match("apps.foo.example.com/prefix", pat_subdomain)
    assert is_registry_match("sub.example.com:80", pat_subdomain)
    assert not is_registry_match("https://example.com:443/prefix", pat_subdomain)
    assert not is_registry_match("docker.io/library/my", pat_subdomain)
    assert not is_registry_match("https://hello.example.bar", pat_subdomain)

    assert is_registry_match("registry:80/prefix", pat_matchport)
    assert is_registry_match("registry/myapp", pat_matchport)
    assert is_registry_match("registry:443/myap", pat_matchport)
    assert not is_registry_match("https://example.com:443/prefix", pat_matchport)
    assert not is_registry_match("docker.io/library/my", pat_matchport)
    assert not is_registry_match("https://hello.registry/myapp", pat_matchport)


if __name__ == '__main__':
    test_is_registry_match()
