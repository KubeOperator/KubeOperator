from rest_framework import serializers


class GitSerializer(serializers.Serializer):
    repo = serializers.URLField()
    branch = serializers.CharField(max_length=64, initial='master', required=False)
