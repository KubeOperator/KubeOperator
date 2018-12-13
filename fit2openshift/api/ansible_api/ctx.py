# -*- coding: utf-8 -*-
#

from functools import partial

from werkzeug.local import Local, LocalProxy


_thread_locals = Local()


def set_current_project(p):
    from .models import Project
    if isinstance(p, str):
        p = Project.objects.get(name=p)
    setattr(_thread_locals, 'current_project', p)


def change_to_root():
    from .models import Project
    set_current_project(Project.root_project())


def get_current_project():
    p = _find('current_project')
    return p


def _find(attr):
    return getattr(_thread_locals, attr, None)


current_project = LocalProxy(partial(_find, 'current_project'))

