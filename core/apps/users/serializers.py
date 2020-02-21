from rest_framework import serializers
from kubeops_api.models import Item
from users.models import User


class ProfileSerializer(serializers.ModelSerializer):
    name = serializers.SerializerMethodField()
    current_item = serializers.SlugRelatedField(
        many=False, required=False, slug_field='name', queryset=Item.objects.all()
    )

    class Meta:
        model = User
        fields = ['id', 'name', 'current_item', 'last_login', 'is_superuser', 'email', 'is_staff', 'is_active',
                  'date_joined', 'items']

    @staticmethod
    def get_name(obj):
        if obj.first_name or obj.last_name:
            return " ".join([obj.first_name, obj.last_name])
        else:
            return obj.username


class UserSerializer(serializers.ModelSerializer):
    current_item = serializers.SlugRelatedField(
        many=False, required=False, slug_field='name', queryset=Item.objects.all()
    )

    class Meta:
        model = User
        fields = [
            'id', 'username', 'email',
            'is_superuser', 'is_active', 'date_joined', 'last_login', 'current_item', 'items'
        ]
        read_only_fields = ['date_joined', 'last_login']


class UserCreateUpdateSerializer(UserSerializer):

    def save(self, **kwargs):
        password = self.validated_data.pop("password", None)
        instance = super().save(**kwargs)
        if password:
            instance.set_password(password)
            instance.save()
        return instance

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        names.append('password')
        return names
