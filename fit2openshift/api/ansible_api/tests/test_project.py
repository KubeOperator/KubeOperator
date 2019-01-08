import json

from .base import BaseTestCase


class ProjectTestCase(BaseTestCase):
    project_url = '/api/v1/projects/'
    project_id = None

    def test_project_create(self):
        data = {
            "name": "Test project",
            "inventory_data": {
                "groups": [
                    {
                        "name": "CentOS",
                        "vars": {
                            "desc": "Centos group"
                        },
                        "children": ["redhat"]
                    },
                    {
                        "name": "redhat",
                    }
                ],
                "hosts": [
                    {
                        "hostname": "192.168.1.1",
                        "ip": "192.168.1.1",
                        "port": 22,
                        "username": "root",
                        "password": "redaht",
                        "groups": ["redhat"]
                    },
                    {
                        "hostname": "gaga",
                        "ip": "192.168.2.1",
                        "port": 23,
                        "username": "admin",
                        "password": "redhat",
                        "groups": ["redhat"],
                    }
                ]
            },
            "options": {"yes": "I do"},
            "comment": "It's a comment"
        }
        url = self.project_url
        response = self.client.post_json(url, data)
        self.assertEqual(response.status_code, 201, "Response code not 201")

        inventory = data.pop('inventory_data', None)
        inventory_resp = response.data.pop('inventory_data', None)

        hosts = {i['hostname']: i for i in inventory.get('hosts')}
        hosts_resp = inventory_resp.get('hosts')
        groups = inventory.get('groups')
        group_resp = inventory_resp.get('groups')

        self.assertEqual(len(hosts or []), len(hosts_resp or []), "Hosts not equal")
        self.assertEqual(len(groups), len(group_resp), "Groups not equal")

        for k, v in response.data.items():
            if k == 'id':
                self.project_id = v
            if k in ('id', 'created_by', 'date_created'):
                continue
            self.assertEqual(v, data.get(k))

    def test_get_project_detail(self):
        self.test_project_create()
        url = '{}{}/'.format(self.project_url, self.project_id)
        print(url)
        resp = self.client.get_json(url)
        print(resp.data)
