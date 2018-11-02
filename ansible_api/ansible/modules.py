# -*- coding: utf-8 -*-
#
import os
import importlib

from django.core.cache import cache
import yaml
import ansible


def generate_ansible_module_index():
    cache_key = "ANSIBLE_MODULES_INDEX"
    modules_cached = cache.get(cache_key)
    if modules_cached:
        return modules_cached

    modules = {}
    module_path = os.path.join(ansible.__path__[0], 'modules')
    for catalog in os.listdir(module_path):
        if not os.path.isdir(os.path.join(module_path, catalog)):
            continue
        modules[catalog] = {}
        catalog_path = os.path.join(module_path, catalog)
        for module_file in os.listdir(catalog_path):
            if module_file == '__init__.py' or not module_file.endswith('.py'):
                continue
            module_name = module_file.replace('.py', '', 1)
            module = importlib.import_module('ansible.modules.{}.{}'.format(catalog, module_name))
            document = yaml.load(getattr(module, 'DOCUMENTATION'))
            modules[catalog][module_name] = document
            del module_name
    cache.set(cache_key, modules, 3600*7)
    return modules


class AnsibleModules:
    index = generate_ansible_module_index()

    @classmethod
    def category(cls):
        return cls.index.keys()

    @classmethod
    def category_with_modules(cls):
        category = {}
        for k, v in cls.index.items():
            if not v:
                continue
            category[k] = v.keys()
        return category

    @classmethod
    def modules(cls):
        modules = {}
        for v in cls.index.values():
            modules.update(v)
        return modules

    @classmethod
    def modules_linux(cls):
        return {k: v for k, v in cls.modules().items() if not k.startswith('win')}

    @classmethod
    def find_module(cls, name):
        for k, v in cls.index.items():
            module = v.get(name)
            if module:
                return module
