from rest_framework import serializers

from ops.models import Script

__all__ = ["ScriptSerializer"]


class ScriptSerializer(serializers.ModelSerializer):
    class Meta:
        model = Script
        fields = ["id", "name", "type", "content", "date_created"]
        read_only_fields = ['id', 'date_created']
