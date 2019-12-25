from common.ssh import SSHClient


def get_gpu_device(ssh_config):
    client = SSHClient(ssh_config)
    cmd = "lspci | grep -i nvidia"
    result, code = client.run_cmd(cmd)
    if code == 0:
        return result
