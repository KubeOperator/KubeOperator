import os
import threading
import uuid
from time import sleep

import yaml
from django.db import models

from ansible_api.models.mixins import AbstractExecutionModel
from cloud_provider import get_cloud_client
from common import models as common_models
from fit2ansible import settings
from django.utils.translation import ugettext_lazy as _
from kubeops_api.models.host import Host


class CloudProviderTemplate(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = common_models.JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    template_dir = os.path.join(settings.BASE_DIR, 'resource', 'clouds')

    @property
    def path(self):
        return os.path.join(self.template_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.template_dir):
            full_path = os.path.join(cls.template_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class Region(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    template = models.ForeignKey('CloudProviderTemplate', on_delete=models.SET_NULL, null=True)
    cloud_region = models.CharField(max_length=128, null=True, default=None)
    vars = common_models.JsonDictTextField(default={})
    comment = models.CharField(max_length=128, blank=True, null=True, verbose_name=_("Comment"))

    @property
    def zone_size(self):
        zones = Zone.objects.filter(region=self)
        return len(zones)

    @property
    def cluster_size(self):
        clusters = []
        plans = Plan.objects.filter(region=self)
        for plan in plans:
            from kubeops_api.models.cluster import Cluster
            cs = Cluster.objects.filter(plan=plan)
            for c in cs:
                clusters.append(c)
        return len(clusters)

    @property
    def image_ovf_path(self):
        return self.vars['image_ovf_path']

    @property
    def image_vmdk_path(self):
        return self.vars['image_vmdk_path']

    @property
    def image_name(self):
        return self.vars['image_name']

    def set_vars(self):
        meta = self.template.meta.get('region', None)
        if meta:
            _vars = meta.get('vars', {})
            self.vars.update(_vars)
            self.save()

    def on_region_create(self):
        self.set_vars()

    def to_dict(self):
        dic = {
            "region": self.cloud_region
        }
        dic.update(self.vars)
        return dic


class Zone(models.Model):
    ZONE_STATUS_READY = "READY"
    ZONE_STATUS_INITIALIZING = "INITIALIZING"
    ZONE_STATUS_ERROR = "ERROR"
    ZONE_STATUS_CHOICES = (
        (ZONE_STATUS_READY, 'READY'),
        (ZONE_STATUS_INITIALIZING, 'INITIALIZING'),
        (ZONE_STATUS_ERROR, 'ERROR'),
    )
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    vars = common_models.JsonDictTextField(default={})
    region = models.ForeignKey('Region', on_delete=models.CASCADE, null=True)
    cloud_zone = models.CharField(max_length=128, null=True, default=None)
    status = models.CharField(max_length=64, choices=ZONE_STATUS_CHOICES, null=True)

    @property
    def cluster_size(self):
        clusters = []
        plans = Plan.objects.filter(zone=self)
        for plan in plans:
            from kubeops_api.models.cluster import Cluster
            cs = Cluster.objects.all().filter(plan=plan)
            for c in cs:
                clusters.append(c)
        return len(clusters)

    @property
    def plan_size(self):
        return len(Plan.objects.filter(zone=self))

    def change_status(self, status):
        self.status = status
        self.save()

    def create_image(self):
        try:
            self.change_status(Zone.ZONE_STATUS_INITIALIZING)
            client = get_cloud_client(self.region.vars)
            client.create_image(zone=self)
            self.change_status(Zone.ZONE_STATUS_READY)
        except Exception as e:
            print(e.args)
            self.change_status(Zone.ZONE_STATUS_ERROR)

    def on_zone_create(self):
        thread = threading.Thread(target=self.create_image)
        thread.start()

    def to_dict(self):
        dic = {
            "name": self.cloud_zone
        }
        dic.update(self.vars)
        return dic

    @property
    def provider(self):
        return self.region.template.name

class Plan(models.Model):
    DEPLOY_TEMPLATE_SINGLE = "SINGLE"
    DEPLOY_TEMPLATE_MULTIPLE = "MULTIPLE"
    DEPLOY_TEMPLATE_CHOICES = (
        (DEPLOY_TEMPLATE_SINGLE, 'single'),
        (DEPLOY_TEMPLATE_MULTIPLE, 'multiple'),
    )

    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    zone = models.ForeignKey('Zone', null=True, on_delete=models.CASCADE)
    region = models.ForeignKey('Region', null=True, on_delete=models.CASCADE)
    zones = models.ManyToManyField('Zone', related_name='zones')
    deploy_template = models.CharField(choices=DEPLOY_TEMPLATE_CHOICES, default=DEPLOY_TEMPLATE_SINGLE, max_length=128)
    vars = common_models.JsonDictTextField(default={})

    @property
    def mixed_vars(self):
        _vars = self.vars.copy()
        _vars.update(self.region.to_dict())
        zones = []
        if self.zones:
            for zone in self.zones.all():
                zones.append(zone.to_dict())
        if self.zone:
            zones.append(self.zone.to_dict())
        _vars['zones'] = zones
        return _vars


    @property
    def compute_models(self):
        return self.region.template.meta["plan"]["models"]


class TerraformHost(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=255, unique=True, verbose_name=_('Name'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    domain = models.CharField(max_length=255)
    cpu = models.IntegerField(default=0)
    memory = models.IntegerField(default=0)
    short_name = models.CharField(max_length=128)
    host_name = models.CharField(max_length=255)
    role = models.CharField(max_length=32)
    ip = models.CharField(max_length=32)
    zone_vars = common_models.JsonDictTextField(default={})
    host = models.ForeignKey('kubeops_api.Host', on_delete=models.CASCADE, null=True)

    def to_dict(self):
        return {
            "name": self.name,
            "domain": self.domain,
            "cpu": self.cpu,
            "memory": self.memory,
            "short_name": self.short_name,
            "host_name": self.host_name,
            "role": self.role,
            "ip": self.ip,
            "zone": self.zone_vars
        }

    def create_host(self):
        username = 'root'
        password = 'KubeOperator@2019'
        host = Host.objects.create(
            name=self.name,
            ip=self.ip,
            username=username,
            password=password
        )
        host.gather_info()
        self.host = host
        self.save()
