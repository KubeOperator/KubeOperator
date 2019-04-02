# -*- coding: utf-8 -*-
#

from .ansible.inventory import BaseInventory


__all__ = [
    'AdHocInventory', 'LocalModelInventory'
]


class AdHocInventory(BaseInventory):
    """
    JMS Inventory is the manager with jumpserver assets, so you can
    write you own manager, construct you inventory_obj
    """
    def __init__(self, assets, nodes=None, run_as_admin=False,
                 run_as=None, become_info=None, vars=None):
        """
        :param assets: ["uuid1", ]
        :param run_as_admin: True 是否使用管理用户去执行, 每台服务器的管理用户可能不同
        :param run_as: 是否统一使用某个系统用户去执行
        :param become_info: 是否become成某个用户去执行
        """
        self.assets = assets or []
        self.nodes = nodes or []
        self.run_as_admin = run_as_admin
        self.run_as = run_as
        self.become_info = become_info
        self.vars = vars or {}
        hosts, groups = self.parse_resource()
        super().__init__(host_list=hosts, group_list=groups)

    def parse_resource(self):
        assets = self.get_all_assets()
        hosts = []
        nodes = set()
        groups = []

        for asset in assets:
            info = self.convert_to_ansible(asset)
            hosts.append(info)
            nodes.update(set(asset.nodes.all()))

        if self.become_info:
            for host in hosts:
                host.update(self.become_info)

        for node in nodes:
            groups.append({
                'name': node.value,
                'children': [n.value for n in node.get_children()]
            })
        return hosts, groups

    def get_all_assets(self):
        assets = set(self.assets)
        for node in self.nodes:
            _assets = set(node.get_all_assets())
            assets.update(_assets)
        return assets

    def convert_to_ansible(self, asset):
        info = {
            'id': asset.id,
            'hostname': asset.hostname,
            'ip': asset.ip,
            'port': asset.port,
            'vars': dict(),
            'groups': [],
        }
        if asset.domain and asset.domain.has_gateway():
            info["vars"].update(self.make_proxy_command(asset))
        if self.run_as_admin:
            info.update(asset.get_auth_info())
        if self.run_as:
            info.update(self.get_run_user_info(asset))
        for node in asset.nodes.all():
            info["groups"].append(node.value)
        for label in asset.labels.all():
            info["vars"].update({
                label.name: label.value
            })
            info["groups"].append("{}:{}".format(label.name, label.value))
        if asset.domain:
            info["vars"].update({
                "domain": asset.domain.name,
            })
            info["groups"].append("domain_"+asset.domain.name)
        for k, v in self.vars.items():
            if not k.startswith('__'):
                info['vars'].update({
                    k: v
                })
        host_vars = self.vars.get("__{}".format(asset.id), {})
        for k, v in host_vars.items():
            info['vars'].update({
                k: v
            })
        return info

    def get_run_user_info(self, asset):
        info = self.run_as.get_auth(asset)._to_secret_json()
        return info

    @staticmethod
    def make_proxy_command(asset):
        gateway = asset.domain.random_gateway()
        proxy_command_list = [
            "ssh", "-p", str(gateway.port),
            "{}@{}".format(gateway.username, gateway.ip),
            "-W", "%h:%p", "-q",
        ]

        if gateway.password:
            proxy_command_list.insert(
                0, "sshpass -p {}".format(gateway.password)
            )
        if gateway.private_key:
            proxy_command_list.append("-i {}".format(gateway.private_key_file))

        proxy_command = "'-o ProxyCommand={}'".format(
            " ".join(proxy_command_list)
        )
        return {"ansible_ssh_common_args": proxy_command}


class JMSInventory(BaseInventory):
    system = object()

    def __init__(self, nodes, vars=None):
        """
        :param nodes: {"node": {"asset": set(user1, user2)}}
        :param vars:
        """
        self.nodes = nodes or {}
        self.vars = vars or {}
        host_list, group_list = self.parse_resource()
        super().__init__(host_list, group_list)

    @classmethod
    def from_json(cls, data):
        host_list = data.get('host_list', [])
        group_list = data.get('group_list', [])
        super().__init__(host_list, group_list)

    def parse_all_hosts(self):
        host_list = []
        for assets in self.nodes.values():
            for asset in assets:
                _vars = {
                    'ansible_ssh_host': asset.ip,
                    'ansible_ssh_port': asset.port,
                }
                if asset.domain and asset.domain.has_gateway():
                    _vars.update(self.make_proxy_command(asset))
                for k, v in self.vars.items():
                    if not k.startswith('__'):
                        _vars.update({
                            k: v
                        })
                host_vars = self.vars.get("__{}".format(asset.hostname), {})
                for k, v in host_vars.items():
                    _vars.update({
                        k: v
                    })
                host_list.append({
                    'hostname': asset.hostname,
                    'vars': _vars
                })
        return host_list

    def parse_label(self):
        pass

    def parse_users(self):
        pass

    def parse_resource(self):
        """
        host_list: [{
            "hostname": "",
            "vars": {},
          },
        ]
        group_list: [{
            "name": "",
            "hosts": ["",],
            "children": ["",],
            "vars": {}
        },]
        :return: host_list, group_list
        """
        host_list = self.parse_all_hosts()
        group_list = []
        return host_list, group_list

    @staticmethod
    def get_run_user_info(user, asset):
        info = user.get_auth(asset)._to_secret_json()
        return info

    @staticmethod
    def make_proxy_command(asset):
        gateway = asset.domain.random_gateway()
        proxy_command_list = [
            "ssh", "-p", str(gateway.port),
            "{}@{}".format(gateway.username, gateway.ip),
            "-W", "%h:%p", "-q",
        ]

        if gateway.password:
            proxy_command_list.insert(
                0, "sshpass -p {}".format(gateway.password)
            )
        if gateway.private_key:
            proxy_command_list.append("-i {}".format(gateway.private_key_file))

        proxy_command = "'-o ProxyCommand={}'".format(
            " ".join(proxy_command_list)
        )
        return {"ansible_ssh_common_args": proxy_command}


class LocalModelInventory(BaseInventory):
    def __init__(self, inventory):
        """
        :param inventory: .models Inventory Object
        """
        self.inventory = inventory
        data = self.parse_resource()
        super().__init__(data)

    def parse_resource(self):
        groups = self._parse_groups()
        hosts = self._parse_hosts()
        return {'hosts': hosts, 'groups': groups}

    def _parse_hosts(self):
        hosts = []
        _hosts = set(self.inventory.hosts)
        for group in self.inventory.groups:
            _hosts.update(set(group.hosts.all()))

        for host in _hosts:
            _vars = host.vars or {}
            _vars.update({
                'ansible_ssh_host': host.ip or host.name,
                'ansible_ssh_port': host.port or 22,
                'ansible_ssh_user': host.username,
                'ansible_ssh_pass': host.password,
            })
            hosts.append({
                'hostname': host.name,
                'vars': _vars
            })
        return hosts

    def _parse_groups(self):
        groups = []
        for group in self.inventory.groups:
            groups.append({
                'name': group.name,
                'vars': group.vars or {},
                'hosts': group.hosts.all().values_list('name', flat=True),
                'children': group.children.all().values_list('name', flat=True)
            })
        return groups


class WithHostInfoInventory(BaseInventory):
    def __init__(self, inventory_data):
        self.inventory_data = inventory_data
        data = self.parse_resource()
        super().__init__(data)

    def parse_resource(self):
        groups = self._parse_groups()
        hosts = self._parse_hosts()
        return {'hosts': hosts, 'groups': groups}

    def _parse_hosts(self):
        hosts = []
        _hosts = self.inventory_data.get('hosts', [])
        for host in _hosts:
            _vars = host.get('vars', {})
            _vars.update({
                'ansible_ssh_host': host.get('ip') or host['name'],
                'ansible_ssh_port': host.get('port') or 22,
                'ansible_ssh_user': host.get('username'),
                'ansible_ssh_pass': host.get('password'),
            })
            hosts.append({
                'hostname': host['name'],
                'vars': _vars
            })
        return hosts

    def _parse_groups(self):
        groups = []
        for group in self.inventory_data.get('groups', []):
            groups.append({
                'name': group.get('name'),
                'vars': group.get('vars', {}),
                'hosts': group.get('hosts', []),
                'children': group.get('children', [])
            })
        return groups
