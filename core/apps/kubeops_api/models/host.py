import logging
import uuid
from django.db import models
from ansible_api.models.inventory import BaseHost
from ansible_api.models.utils import name_validator
from common.ssh import SshConfig
from kubeoperator.settings import NODE_CREDENTIAL
from kubeops_api.adhoc import gather_host_info
from kubeops_api.models.credential import Credential
from common import models as common_models
from kubeops_api.utils.gpu import get_gpu_device

__all__ = ['Host', 'Volume', 'GPU']
logger = logging.getLogger('kubeops')


class Host(BaseHost):
    HOST_STATUS_RUNNING = "RUNNING"
    HOST_STATUS_CREATING = "CREATING"
    HOST_STATUS_UNKNOWN = "UNKNOWN"
    HOST_STATUS_UPDATING = "UPDATING"
    DEPLOY_TEMPLATE_CHOICES = (
        (HOST_STATUS_RUNNING, 'running'),
        (HOST_STATUS_CREATING, 'creating'),
        (HOST_STATUS_UNKNOWN, "unknown"),
        (HOST_STATUS_UPDATING, "updating")
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    node = models.ForeignKey('Node', default=None, null=True, related_name='node', on_delete=models.SET_NULL)
    name = models.CharField(max_length=128, validators=[name_validator], unique=True)
    credential = models.ForeignKey("kubeops_api.Credential", null=True, on_delete=models.SET_NULL)
    memory = models.fields.BigIntegerField(default=0)
    os = models.fields.CharField(max_length=128, default="")
    os_version = models.fields.CharField(max_length=128, default="")
    cpu_core = models.fields.IntegerField(default=0)
    volumes = models.ManyToManyField('Volume')
    gpus = models.ManyToManyField('GPU')
    zone = models.ForeignKey('cloud_provider.Zone', null=True, on_delete=models.CASCADE)
    status = models.CharField(choices=DEPLOY_TEMPLATE_CHOICES, default=HOST_STATUS_UNKNOWN, max_length=128)
    username = models.CharField(max_length=256, default=NODE_CREDENTIAL['username'])
    password = common_models.EncryptCharField(max_length=4096, blank=True, null=True,
                                              default=NODE_CREDENTIAL['password'])
    auto_gather_info = models.BooleanField(default=True, null=True)
    conditions = models.ManyToManyField("Condition")
    gpu = models.BooleanField(default=False)
    gpu_info = models.CharField(default="", max_length=128)
    gpu_num = models.IntegerField(default=0)

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
    def has_gpu(self):
        gpus = self.gpus.all()
        if gpus and len(gpus) > 0:
            return True
        else:
            return False

    @property
    def region(self):
        if self.zone:
            return self.zone.region.name

    def delete(self, using=None, keep_parents=False):
        if self.zone:
            self.zone.recover_ip(self.ip)
        super().delete(using=None, keep_parents=False)

    def to_ssh_config(self):
        return SshConfig(
            self.ip,
            self.port,
            self.username,
            self.password,
            self.private_key,
        )

    def health_check(self):
        from kubeops_api.models.health.host_health import HostHealthCheck
        health_check = HostHealthCheck(host=self)
        health_check.run()

    def gather_info(self, retry=1):
        self.host.status = Host.HOST_STATUS_UPDATING
        try:
            logger.info("host: {}  gather host info ".format(self.name))
            facts = gather_host_info(ip=self.ip, port=self.port, username=self.username, retry=retry,
                                     password=self.password,
                                     private_key_path=self.private_key_path)
            self.memory = facts["ansible_memtotal_mb"]
            cpu_cores = facts["ansible_processor_cores"]
            cpu_count = facts["ansible_processor_count"]
            vcpus = facts.get("ansible_processor_vcpus", None)
            if vcpus:
                self.cpu_core = int(vcpus)
            else:
                self.cpu_core = int(cpu_cores) * int(cpu_count)
            self.os = facts["ansible_distribution"]
            self.os_version = facts["ansible_distribution_version"]
            devices = facts["ansible_devices"]
            volumes = []
            for name in devices:
                if not name.startswith(('dm', 'loop', 'sr')):
                    volume = Volume(name='/dev/' + name)
                    volume.size = devices[name]['size']
                    volume.save()
                    volumes.append(volume)
            self.volumes.set(volumes)
            msg = get_gpu_device(self.to_ssh_config())
            gpus = []
            if msg:
                host_gpus = str(msg).split('\n')
                for hg in host_gpus:
                    g = GPU()
                    g.name = hg[hg.index("[") + 1:hg.index("]")]
                    g.save()
                    gpus.append(g)
            self.gpus.set(gpus)
            self.status = Host.HOST_STATUS_RUNNING
            self.save()
        except Exception as e:
            self.status = Host.HOST_STATUS_UNKNOWN
            self.cpu_core = 0
            self.memory = 0
            self.os = "未知"
            self.volumes.set([])
            self.os_version = ""
            self.save()
            logger.error("host: {}  gather host info".format(self.name), exc_info=True)

    class Meta:
        ordering = ('name',)


class Volume(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128)
    size = models.CharField(max_length=16)

    class Meta:
        ordering = ('size',)


class GPU(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=256)
