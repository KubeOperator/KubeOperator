# -*- coding: utf-8 -*-
#

from rest_framework import serializers
from django.contrib.auth import get_user_model


class ProfileSerializer(serializers.ModelSerializer):
    name = serializers.SerializerMethodField()

    class Meta:
        model = get_user_model()
        exclude = [
            'password', 'first_name', 'last_name',
        ]

    @staticmethod
    def get_name(obj):
        if obj.first_name or obj.last_name:
            return " ".join([obj.first_name, obj.last_name])
        else:
            return obj.username


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = get_user_model()
        fields = [
            'id', 'username', 'email',
            'is_superuser', 'is_active', 'date_joined', 'last_login'
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
