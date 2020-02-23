from django.contrib.auth.models import User
from rest_framework import serializers
from kubeops_api.models import Item
from users.models import Profile

__all__ = ["UserSerializer", "ProfileSerializer", "UserCreateUpdateSerializer"]


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = [
            'id', 'username', 'email',
            'is_superuser', 'is_active', 'date_joined', 'last_login',
        ]
        read_only_fields = ['date_joined', 'last_login']


class UserCreateUpdateSerializer(UserSerializer):

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

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        return names


class ProfileSerializer(serializers.ModelSerializer):
    current_item = serializers.SlugRelatedField(required=False, slug_field='name', queryset=Item.objects.all())
    user = UserSerializer(read_only=True)

    class Meta:
        model = Profile
        fields = ["current_item", "user"]
        read_only_fields = ['user']
