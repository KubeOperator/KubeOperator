from rest_framework import serializers
from kubeops_api.models.item import Item
from users.models import User
from users.serializers import UserSerializer


class ItemSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name', 'description', 'date_created']


class ItemUserReadSerializer(serializers.ModelSerializer):
    users = UserSerializer(many=True)

    class Meta:
        model = Item
        fields = ['id', 'name', 'users']


class ItemUserSerializer(serializers.ModelSerializer):
    users = serializers.SlugRelatedField(
        many=True,
        queryset=User.objects.all(),
        required=True,
        slug_field='id'
    )

    class Meta:
        model = Item
        fields = ['users']
