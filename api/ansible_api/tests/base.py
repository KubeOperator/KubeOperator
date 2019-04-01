# -*- coding: utf-8 -*-
#
import json

from django.test import TestCase, Client
from django.contrib.auth.models import User

BOUNDARY = 'BoUnDaRyStRiNg'
MULTIPART_CONTENT = 'multipart/form-data; boundary=%s' % BOUNDARY

content_type_json = 'application/json'


class BaseClient(Client):

    def post_json(self, path, data, **kwargs):
        return super().post(path, json.dumps(data),
                            content_type=content_type_json, **kwargs)

    def get_json(self, path, **kwargs):
        return super().get(path, content_type=content_type_json, **kwargs)


class BaseTestCase(TestCase):
    def setUp(self):
        admin = User.objects.create(
            username='admin', is_superuser=True, is_active=True, is_staff=True
        )
        admin.set_password('redhat123')
        admin.save()
        self.client = BaseClient()
        self.client.login(username='admin', password='redhat123')
