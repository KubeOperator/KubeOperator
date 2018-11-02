# -*- coding: utf-8 -*-
#
import uuid

from django.db import models

from ..ctx import current_project, set_current_project


class ProjectResourceManager(models.Manager):
    def get_queryset(self):
        queryset = super(ProjectResourceManager, self).get_queryset()
        if not current_project:
            return queryset
        if current_project.is_real():
            queryset = queryset.filter(project=current_project.id)
        return queryset

    def create(self, **kwargs):
        if 'project' not in kwargs and current_project.is_real():
            kwargs['project'] = current_project._get_current_object()
        return super().create(**kwargs)

    def all(self):
        if current_project:
            return super().all()
        else:
            return self

    def set_current_org(self, project):
        set_current_project(project)
        return self


class AbstractProjectResourceModel(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    project = models.ForeignKey('Project', on_delete=models.CASCADE)
    objects = ProjectResourceManager()

    name = 'Not-Sure'

    class Meta:
        abstract = True

    def __str__(self):
        return '{}: {}'.format(self.project, self.name)

    def save(self, force_insert=False, force_update=False, using=None,
             update_fields=None):
        if not hasattr(self, 'project') and current_project.is_real():
            self.project = current_project._get_current_object()
        return super().save(force_insert=force_insert, force_update=force_update,
                            using=using, update_fields=update_fields)
