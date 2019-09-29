import uuid
from django.db import models
from ansible_api.models.inventory import BaseHost
from ansible_api.models.utils import name_validator
from kubeops_api.adhoc import gather_host_info
from kubeops_api.models.credential import Credential

__all__ = ['Host']


class Host(BaseHost):
    HOST_STATUS_RUNNING = "RUNNING"
    HOST_STATUS_CREATING = "CREATING"
    HOST_STATUS_UNKNOWN = "UNKNOWN"
    DEPLOY_TEMPLATE_CHOICES = (
        (HOST_STATUS_RUNNING, 'running'),
        (HOST_STATUS_CREATING, 'creating'),
        (HOST_STATUS_UNKNOWN, "unknown")
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    node = models.ForeignKey('Node', default=None, null=True, related_name='node',
                             on_delete=models.SET_NULL)
    name = models.CharField(max_length=128, validators=[name_validator], unique=True)
    credential = models.ForeignKey("kubeops_api.Credential", null=True, on_delete=models.SET_NULL)
    memory = models.fields.BigIntegerField(default=0)
    os = models.fields.CharField(max_length=128, default="")
    os_version = models.fields.CharField(max_length=128, default="")
    cpu_core = models.fields.IntegerField(default=0)
    volumes = models.ManyToManyField('Volume')
    zone = models.ForeignKey('cloud_provider.Zone', null=True, on_delete=models.CASCADE)
    status = models.CharField(choices=DEPLOY_TEMPLATE_CHOICES, default=HOST_STATUS_UNKNOWN, max_length=128)

    def full_host_credential(self):
        if self.credential:
            self.username = self.credential.username
            if self.credential.type == Credential.CREDENTIAL_TYPE_PASSWORD:
                self.password = self.credential.password
            else:
                self.private_key = self.credential.private_key
            self.save()

    @property
    def cluster(self):
        if self.node:
            return self.node.project.name

    @property
    def region(self):
        if self.zone:
            return self.zone.region.name

    def gather_info(self):
        facts = gather_host_info(self.ip, self.username, self.password)
        self.memory = facts["ansible_memtotal_mb"]
        cpu_cores = facts["ansible_processor_cores"]
        cpu_count = facts["ansible_processor_count"]
        self.cpu_core = int(cpu_cores) * int(cpu_count)
        self.os = facts["ansible_distribution"]
        self.os_version = facts["ansible_distribution_version"]
        self.save()
        devices = facts["ansible_devices"]
        volumes = []
        for name in devices:
            if not name.startswith(('dm', 'loop', 'sr')):
                volume = Volume(name='/dev/' + name)
                volume.size = devices[name]['size']
                volume.save()
                volumes.append(volume)
        self.volumes.set(volumes)
        self.status = Host.HOST_STATUS_RUNNING
        self.save()

    class Meta:
        ordering = ('name',)


class Volume(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128)
    size = models.CharField(max_length=16)

    class Meta:
        ordering = ('size',)
