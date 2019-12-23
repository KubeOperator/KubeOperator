from common.ssh import SSHClient


def test_ssh(ssh_config):
    result = True
    client = SSHClient(ssh_config)
    try:
        client.ping()
    except Exception as e:
        result = False
        return result, e
    return result, None
