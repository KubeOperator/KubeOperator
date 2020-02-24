from ansible_api.serializers.inventory import HostReadSerializer
from rest_framework import serializers
from cloud_provider.models import Zone
from common.ssh import SshConfig, SSHClient
from kubeops_api.models import Credential, Condition
from kubeops_api.models.host import Volume, GPU, Host

__all__ = ["VolumeSerializer", "GPUSerializer", "ConditionSerializer", "HostSerializer"]


class HostConnectedException(Exception):
    pass


def is_ip_exists(ip):
    hosts = Host.objects.filter(ip=ip)
    return len(hosts) > 0


def is_host_connected(ip, port, credential):
    config = SshConfig(host=ip, username=credential.username, password=credential.password,
                       private_key=credential.private_key,
                       port=port)
    client = SSHClient(config)
    return client.ping()


class HostSerializerMixin(serializers.ModelSerializer):

    def validate(self, attrs):
        ip = attrs.get("ip")
        credential = attrs.get("credential")
        port = attrs.get("port", 22)
        if is_ip_exists(ip):
            raise serializers.ValidationError("ip {} already exists!".format(ip))
        if not is_host_connected(ip, port, credential):
            raise serializers.ValidationError("can not connected host: {} with given credential".format(ip))
        return attrs

    def save(self, **kwargs):
        self.instance = super().save(**kwargs)
        self.instance.gather_info()
        return self.instance


class VolumeSerializer(serializers.ModelSerializer):
    class Meta:
        model = Volume
        fields = [
            'id', 'name', 'size',
        ]
        read_only_fields = ['id', 'name', 'size', ]


class GPUSerializer(serializers.ModelSerializer):
    class Meta:
        model = GPU
        fields = [
            'id', 'name',
        ]
        read_only_fields = ['id', 'name', ]


class ConditionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Condition
        fields = [
            "type", "status", "message", "reason", "last_time"
        ]
        read_only_fields = ["type", "status", "message", "reason", "last_time"]


class HostSerializer(HostReadSerializer, HostSerializerMixin):
    credential = serializers.SlugRelatedField(
        queryset=Credential.objects.all(),
        slug_field='name', required=False
    )
    zone = serializers.SlugRelatedField(
        queryset=Zone.objects.all(),
        slug_field='name', required=False
    )
    volumes = VolumeSerializer(required=False, many=True)
    gpus = GPUSerializer(required=False, many=True)
    conditions = ConditionSerializer(required=False, many=True)

    class Meta:
        model = Host
        extra_kwargs = HostReadSerializer.Meta.extra_kwargs
        fields = [
            'id', 'name', 'ip', 'port', 'cluster', 'credential', 'memory', 'os', 'os_version', 'cpu_core', 'volumes',
            'zone',
            'region', 'status', 'conditions', 'gpus', "has_gpu"
        ]
        read_only_fields = ['id', 'comment', 'memory', 'os', 'os_version', 'cpu_core', 'volumes', 'zone', 'region',
                            'status', "conditions", 'gpus', "has_gpu"]
