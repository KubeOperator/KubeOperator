"""
Ansible module for determining information about the docker host.

While there are several ansible modules that make use of the docker
api to expose container and image facts in a remote host, they
are unable to return specific information about the host machine
itself. This module exposes the same information obtained through
executing the `docker info` command on a docker host, in json format.
"""

from ansible.module_utils.docker_common import AnsibleDockerClient


def main():
    """Entrypoint for running an Ansible module."""
    client = AnsibleDockerClient()

    client.module.exit_json(
        info=client.info(),
    )


if __name__ == '__main__':
    main()
