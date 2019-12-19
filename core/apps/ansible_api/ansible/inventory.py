# ~*~ coding: utf-8 ~*~
from ansible.inventory.host import Host
from ansible.vars.manager import VariableManager
from ansible.inventory.manager import InventoryManager
from ansible.parsing.dataloader import DataLoader


__all__ = [
    'BaseHost', 'BaseInventory'
]


class BaseHost(Host):
    def __init__(self, host_data):
        """
        初始化
        :param host_data:  {
            "hostname": "",
            "vars": {},
        }
        """
        self.host_data = host_data
        hostname = host_data.get('hostname')
        super().__init__(hostname)
        self.__set_variables()

    def __set_variables(self):
        _vars = self.host_data.get("vars", {})
        for k, v in _vars.items():
            self.set_variable(k, v)

    def __repr__(self):
        return self.name


class BaseInventory(InventoryManager):
    """
    提供生成Ansible inventory对象的方法
    """
    loader_class = DataLoader
    variable_manager_class = VariableManager
    host_manager_class = BaseHost

    def __init__(self, data):
        """
        用于生成动态构建Ansible Inventory. super().__init__ 会自动调用
        host_list 将会parse到all组中

        :param data:
        data = {
            hosts: [{
                "hostname": "",
                "vars": {},
            },
            ...
            ],
            groups: [{
                "name": "",
                "hosts": ["",],
                "children": ["",],
                "vars": {}
            },
            ...
            ]
        }
        """
        self.host_list = data.get('hosts', [])
        self.group_list = data.get('groups',  [])
        assert isinstance(self.host_list, list)
        self.loader = self.loader_class()
        self.variable_manager = self.variable_manager_class()
        super().__init__(self.loader)

    def get_groups(self):
        return self._inventory.groups

    def get_group(self, name):
        return self._inventory.groups.get(name, None)

    def get_host(self, name):
        return self._inventory.hosts.get(name, None)

    def get_or_create_group(self, name):
        group = self.get_group(name)
        if not group:
            self.add_group(name)
            return self.get_or_create_group(name)
        else:
            return group

    def __parse_groups(self):
        for g in self.group_list:
            name = g.get("name")
            group = self.get_or_create_group(name)
            # 添加Host到组
            hosts = g.get("hosts", [])
            for hostname in hosts:
                host = self.get_host(hostname)
                if not host:
                    continue
                group.add_host(host)
            # 添加Children到组
            children = [
                self.get_or_create_group(n) for n in g.get('children', [])
            ]
            for child in children:
                group.add_child_group(child)
            # 组变量
            for k, v in g.get("vars", {}).items():
                group.set_variable(k, v)

    def __parse_hosts(self):
        group_all = self.get_or_create_group('all')
        for host_data in self.host_list:
            host = self.host_manager_class(host_data=host_data)
            self.hosts[host_data['hostname']] = host
            group_all.add_host(host)

    def parse_sources(self, cache=False):
        self.__parse_hosts()
        self.__parse_groups()

    def get_matched_hosts(self, pattern):
        return self.get_hosts(pattern)


