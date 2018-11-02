# -*- coding: utf-8 -*-
#
from django.db import transaction
from django.shortcuts import get_object_or_404

from ..ctx import set_current_project
from ..models import Project


class ProjectObjectMixin:
    is_project_request = False
    project = None
    project_name = ''

    @transaction.atomic
    def dispatch(self, request, *args, **kwargs):
        if kwargs.get('project_name'):
            self.is_project_request = True
            self.project_name = kwargs['project_name']
            self.project = self.get_project()
            set_current_project(self.project)
        return super().dispatch(request, *args, **kwargs)

    def get_project(self):
        if self.project is not None:
            return self.project
        if self.is_project_request:
            self.project = get_object_or_404(Project, name=self.project_name)
        return self.project

    def get_context_data(self, **kwargs):
        kwargs.update({
            'project': self.project
        })
        return super().get_context_data(**kwargs)

    def get_serializer_class(self):
        if hasattr(self, 'action') and self.action in ('list', 'retrieve') \
                and hasattr(self, 'read_serializer_class'):
            print(self.read_serializer_class)
            return self.read_serializer_class
        return super().get_serializer_class()
