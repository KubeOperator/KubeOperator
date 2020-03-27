import logging
import os
import threading
import uuid
from ipaddress import ip_address, ip_interface
import yaml
from django.db import models

from ansible_api.models.mixins import AbstractExecutionModel
from cloud_provider import get_cloud_client
from common import models as common_models
from kubeoperator import settings
from django.utils.translation import ugettext_lazy as _
from kubeops_api.models.host import Host

logger = logging.getLogger('cloud_provider')


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
    ip_used = common_models.JsonListTextField(null=True, default=[])
    status = models.CharField(max_length=64, choices=ZONE_STATUS_CHOICES, null=True)

    @property
    def host_size(self):
        hosts = Host.objects.filter(zone=self)
        return len(hosts)

    def change_status(self, status):
        self.status = status
        self.save()

    def create_image(self):
        try:
            logger.info('upload os image')
            self.change_status(Zone.ZONE_STATUS_INITIALIZING)
            client = get_cloud_client(self.region.vars)
            client.create_image(zone=self)
            self.change_status(Zone.ZONE_STATUS_READY)
        except Exception as e:
            logger.error(msg='upload os image error!', exc_info=True)
            self.change_status(Zone.ZONE_STATUS_ERROR)

    def on_zone_create(self):
        thread = threading.Thread(target=self.create_image)
        thread.start()

    def allocate_ip(self):
        ip = self.ip_pools().pop()
        self.ip_used.append(ip)
        self.save()
        return ip

    def recover_ip(self, ip):
        self.ip_used.remove(ip)
        self.save()

    def to_dict(self):
        dic = {
            "key": "z" + str(self.id).split("-")[3],
            "name": self.cloud_zone,
            "zone_name": self.name,
            "ip_pool": self.ip_pools()
        }
        dic.update(self.vars)
        ip_start = ip_address(self.vars['ip_start'])
        net_mask = self.vars.get('net_mask', None)
        if net_mask:
            interface = ip_interface("{}/{}".format(str(ip_start), net_mask))
            dic["net_mask"] = interface.network.prefixlen
        else:
            dic["net_mask"] = 24
        return dic

    def ip_pools(self):
        ip_pool = []
        ip_start = ip_address(self.vars['ip_start'])
        ip_end = ip_address(self.vars['ip_end'])

        if self.region.template.name == 'openstack':
            while ip_start <= ip_end:
                ip_pool.append(str(ip_start))
                ip_start += 1
            for ip in self.ip_used:
                if ip in ip_pool:
                    ip_pool.remove(ip)
            return ip_pool

        net_mask = self.vars['net_mask']
        interface = ip_interface("{}/{}".format(str(ip_start), net_mask))
        network = interface.network
        for host in network.hosts():
            if ip_start <= host <= ip_end:
                ip_pool.append(str(host))
        for ip in self.ip_used:
            if ip in ip_pool:
                ip_pool.remove(ip)
        return ip_pool

    def ip_available_size(self):
        return len(self.ip_pools())

    def has_plan(self):
        for plan in Plan.objects.all():
            for zone in plan.get_zones():
                if zone.name == self.name:
                    return True
        return False

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
    def provider(self):
        return self.region.vars['provider']

    @property
    def mixed_vars(self):
        _vars = self.vars.copy()
        _vars.update(self.region.to_dict())
        zones = self.get_zones()
        zone_dicts = []
        for zone in zones:
            zone_dicts.append(zone.to_dict())
        _vars['zones'] = zone_dicts
        return _vars

    def get_zones(self):
        zones = []
        if self.zone:
            zones.append(self.zone)
        if self.zones:
            zones.extend(self.zones.all())
        return zones

    def count_ip_available(self):
        zones = self.get_zones()
        num = 0
        for zone in zones:
            num += zone.ip_available_size()
        return num

    @property
    def compute_models(self):
        return {
            "master": self.vars.get('master_model', None),
            "worker": self.vars.get('worker_model', None)
        }
