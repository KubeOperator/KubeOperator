from rest_framework import serializers
from kubeops_api.models.item import Item, ItemRole
from users.models import Profile
from users.serializers import UserSerializer


class ItemSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name', 'description', 'date_created']


class ItemRoleSerializer(serializers.ModelSerializer):
    class Meta:
        model = ItemRole
        fields = ["name"]


class ItemUserReadSerializer(serializers.ModelSerializer):
    users = UserSerializer(many=True)

    class Meta:
        model = Item
        fields = ['id', 'name', 'users']


class ItemUserSerializer(serializers.ModelSerializer):
    users = serializers.SlugRelatedField(
        many=True,
        queryset=Profile.objects.all(),
        required=True,
        slug_field='id'
    )

    class Meta:
        model = Item
        fields = ['users']
