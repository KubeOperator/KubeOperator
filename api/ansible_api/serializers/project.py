# -*- coding: utf-8 -*-
#

import re

from rest_framework import serializers
from django.core.validators import RegexValidator

from ..models import Project


__all__ = [
     "ProjectSerializer"
]


class ProjectSerializer(serializers.ModelSerializer):
    options = serializers.DictField(required=False, default={})
    validated_options = ('forks', 'timeout')

    class Meta:
        model = Project
        fields = [
            'id', 'name',  'options', 'comment',
            'created_by', 'date_created'
        ]
        read_only_fields = ('id', 'created_by', 'date_created')

    def validate_options(self, values):
        for k in values:
            if k not in self.validated_options:
                raise serializers.ValidationError(
                    "Option {} not in {}".format(k, self.validated_options)
                )
        return values

    # 因为drf slug field 存在问题所以，重写了
    # 见 https://github.com/encode/django-rest-framework/pull/6167/
    def get_fields(self):
        fields = super().get_fields()
        name_field = fields.get('name')
        for validator in name_field.validators:
            if isinstance(validator, RegexValidator):
                if validator.regex.pattern == r'^[-a-zA-Z0-9_]+$':
                    validator.regex = re.compile(r'^[-\w]+\Z')
        return fields
