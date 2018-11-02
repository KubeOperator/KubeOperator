# -*- coding: utf-8 -*-
#

from rest_framework import serializers

from ..models import Project
from ..ctx import set_current_project, get_current_project


class ReadSerializerMixin(serializers.Serializer):
    project = serializers.SlugRelatedField(queryset=Project.objects.all(), slug_field='name')


class ProjectSerializerMixin(serializers.Serializer):
    project = serializers.HiddenField(default=get_current_project)

    # Todo: If inherit from front class, user pass project will be using error
