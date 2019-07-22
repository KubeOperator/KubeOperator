import uuid
from django.db import models
from ansible_api.models.inventory import BaseHost
from ansible_api.models.utils import name_validator
from kubeops_api.adhoc import gather_host_info
from kubeops_api.models.credential import Credential

__all__ = ['Host', 'HostInfo']


class Host(BaseHost):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    node = models.ForeignKey('Node', default=None, null=True, related_name='node',
                             on_delete=models.SET_NULL)
    name = models.CharField(max_length=128, validators=[name_validator], unique=True)
    credential = models.ForeignKey("kubeops_api.Credential", null=True, on_delete=models.SET_NULL)

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
        if self.node is not None:
            return self.node.project.name
        else:
            return '无'

    @property
    def info(self):
        return self.infos.all().latest()

    def gather_info(self):
        info = HostInfo.objects.create(host_id=self.id)
        info.gather_info()

    class Meta:
        ordering = ('name',)


class HostInfo(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    memory = models.fields.BigIntegerField(default=0)
    os = models.fields.CharField(max_length=128, default="")
    os_version = models.fields.CharField(max_length=128, default="")
    cpu_core = models.fields.IntegerField(default=0)
    host = models.ForeignKey('Host', on_delete=models.CASCADE, null=True, related_name='infos')
    volumes = models.ManyToManyField('Volume')
    date_created = models.DateTimeField(auto_now_add=True)

    class Meta:
        get_latest_by = 'date_created'

    def gather_info(self):
        facts = gather_host_info(self.host)
        self.memory = facts["ansible_memtotal_mb"]
        self.cpu_core = facts["ansible_processor_cores"]
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


class Volume(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128)
    size = models.CharField(max_length=16)

    class Meta:
        ordering = ('size',)
