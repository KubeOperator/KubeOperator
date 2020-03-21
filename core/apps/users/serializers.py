from django.contrib.auth.models import User
from rest_framework import serializers
from rest_framework.exceptions import ValidationError

from kubeops_api.models import Item
from kubeops_api.models.item import ItemRoleMapping
from users.models import Profile

__all__ = ["UserSerializer", "ProfileSerializer", "UserCreateUpdateSerializer", "ChangeUserPasswordSerializer"]


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = [
            'id', 'username', 'email',
            'is_superuser', 'is_active', 'date_joined', 'last_login',
        ]
        read_only_fields = ['date_joined', 'last_login']


class UserCreateUpdateSerializer(UserSerializer):
    password = serializers.CharField(required=True)

    def create(self, validated_data):
        instance = super().create(validated_data)
        Profile.objects.create(user=instance)
        return instance

    def save(self, **kwargs):
        password = self.validated_data.pop("password", None)
        instance = super().save(**kwargs)
        if password:
            instance.set_password(password)
            instance.save()
        return instance

    class Meta:
        model = User
        fields = [
            'username', 'email', 'is_superuser', 'is_active', 'password'
        ]

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        return names


class ItemReadSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name']


class ItemRoleMappingReadSerializer(serializers.ModelSerializer):
    class Meta:
        model = ItemRoleMapping
        fields = ['role', 'item_name']


class ProfileSerializer(serializers.ModelSerializer):
    user = UserSerializer(read_only=True)
    items = ItemReadSerializer(many=True, read_only=True)
    item_role_mappings = ItemRoleMappingReadSerializer(many=True, read_only=True, default={})

    class Meta:
        model = Profile
        fields = ["id", "user", "items", "item_role_mappings", "source"]
        read_only_fields = ['user']


class ChangeUserPasswordSerializer(serializers.ModelSerializer):
    password = serializers.CharField(required=True)
    original = serializers.CharField(required=True)

    def update(self, instance, validated_data):
        password = validated_data.pop('password')
        original = validated_data.pop('original')
        if instance.check_password(original):
            instance.set_password(password)
            return instance.save()
        else:
            raise ValidationError("original password error")

    class Meta:
        model = User
        fields = ["username", "password", "original"]
        read_only_fields = ['username']
