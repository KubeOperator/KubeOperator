import datetime
import time

import paramiko
from paramiko import SSHException


class SshConfig:
    def __init__(self, host, port, username, password, timeout):
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.timeout = timeout


class SSHClient():
    def __init__(self, config):
        self.config = config

    def run_cmd(self, cmd):
        try:
            client = paramiko.SSHClient()
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            client.connect(
                self.config.host,
                self.config.port,
                self.config.username,
                self.config.password
            )
            session = client.get_transport().open_session()
            session.exec_command(cmd)
            while not session.recv_ready():
                time.sleep(0.5)
            exit_code = session.recv_exit_status()
            if exit_code == 0:
                result = session.recv(4096 * 1024).decode(encoding='UTF-8', errors='strict')
            else:
                result = session.recv_stderr(4096 * 1024).decode(encoding='UTF-8', errors='strict')
            return result, exit_code
        except SSHException as e:
            print("[%s] %s target failed, the reason is %s" % (datetime.datetime.now(), self.config.host, str(e)))

    def ping(self):
        out, code = self.run_cmd("pwd")
        return code == 0
