import datetime
import os
import time

import paramiko
from paramiko import SSHException

from common.utils import ssh_key_string_to_obj
from kubeoperator import settings
from hashlib import md5


class SshConfig:
    def __init__(self, host, port, username, password, private_key=None):
        self.host = host
        self.port = port or 22
        self.username = username
        self.password = password
        self.private_key = private_key


class SSHClient:
    def __init__(self, config):
        self.config = config

    def run_cmd(self, cmd):
        client = paramiko.SSHClient()
        try:
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            client.connect(
                self.config.host,
                self.config.port,
                self.config.username,
                self.config.password,
                key_filename=create_ssh_key(self.config.private_key)
            )
            session = client.get_transport().open_session()
            session.exec_command(cmd)
            while not session.exit_status_ready():
                time.sleep(0.5)
            exit_code = session.recv_exit_status()
            if exit_code == 0:
                result = session.recv(4096 * 1024).decode(encoding='UTF-8', errors='strict').rstrip("\n")
            else:
                result = session.recv_stderr(4096 * 1024).decode(encoding='UTF-8', errors='strict').rstrip("\n")
            return result, exit_code
        except SSHException as e:
            print("[%s] %s target failed, the reason is %s" % (datetime.datetime.now(), self.config.host, str(e)))
        finally:
            client.close()

    def ping(self):
        try:
            out, code = self.run_cmd("pwd")
        except Exception:
            return False
        return code == 0


def create_ssh_key(key):
    private_key_obj = ssh_key_string_to_obj(key, None)
    if not private_key_obj:
        return None
    tmp_dir = os.path.join(settings.BASE_DIR, 'data', 'tmp')
    if not os.path.isdir(tmp_dir):
        os.makedirs(tmp_dir)
    key_name = '.' + md5(key.encode('utf-8')).hexdigest()
    key_path = os.path.join(tmp_dir, key_name)
    if not os.path.exists(key_path):
        private_key_obj.write_private_key_file(key_path)
        os.chmod(key_path, 0o400)
    return key_path
