from rest_framework import serializers

from kubeops_api.models.item import Item, ItemRoleMapping
from users.models import Profile
from users.serializers import ProfileSerializer


class ItemSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name', 'description', 'date_created']


class ItemRoleSerializer(serializers.ModelSerializer):
    class Meta:
        model = ItemRoleMapping
        fields = ["name"]


class ItemUserReadSerializer(serializers.ModelSerializer):
    profiles = ProfileSerializer(many=True, required=True)

    class Meta:
        model = Item
        fields = ['name', 'profiles']


class ItemUserSerializer(serializers.ModelSerializer):
    profiles = serializers.SlugRelatedField(queryset=Profile.objects.all(), many=True, slug_field="id", required=True)
    role_map = serializers.DictField(default={}, required=False)

    def update(self, instance, validated_data):
        role_map = validated_data.pop('role_map')
        profiles = validated_data.get('profiles')
        for p in profiles:
            role = ItemRoleMapping.ITEM_ROLE_VIEWER
            if str(p.id) in role_map:
                role = role_map[str(p.id)]
            defaults = {"item": instance, "role": role, "profile": p}
            ItemRoleMapping.objects.update_or_create(defaults, item=instance, profile=p)
        return super().update(instance, validated_data)

    class Meta:
        model = Item
        fields = ['profiles', 'role_map']
