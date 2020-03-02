from common.ssh import SSHClient, SshConfig


def test_ssh(ssh_config):
    result = True
    client = SSHClient(ssh_config)
    try:
        client.ping()
    except Exception as e:
        result = False
        return result, e
    return result, None


def parse_host_to_ssh_config(host):
    return SshConfig(
        host.ip,
        host.port,
        host.username,
        host.password,
        host.private_key
    )
