from rest_framework import serializers
from kubeops_api.models.item import Item


class ItemSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name', 'description', 'date_created']