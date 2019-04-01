from rest_framework import serializers


class DictField(serializers.DictField):
    def to_representation(self, value):
        if not value or not isinstance(value, dict):
            value = {}
        return super().to_representation(value)


class OutputSerializer(serializers.Serializer):
    pass


class TaskSerializer(serializers.Serializer):
    task = serializers.CharField()
