# -*- coding: utf-8 -*-
#
from collections import defaultdict

from django.core.validators import RegexValidator
from django.utils.translation import ugettext_lazy as _


name_validator = RegexValidator(regex=r'^[a-zA-Z0-9_\-\.]+$', message=_(
    'Enter a valid name consisting of Unicode letters, '
    'numbers, underscores, or hyphens, or dot'
))


def format_result_as_list(result):
    _result = defaultdict(list)
    _result['success'] = result.pop('success', False)

    for status, res in result.items():
        if not isinstance(res, dict):
            _result[status] = res
            continue
        for hostname, _tasks in res.items():
            tasks = []
            for task_name, detail in _tasks.items():
                detail["task"] = task_name
                tasks.append(detail)
            _result[status].append({"hostname": hostname, "tasks": tasks})
    return _result


def format_results_as_list(results):
    _raw = results.get('raw', {})
    _summary = results.get("summary", {})
    results['raw'] = format_result_as_list(_raw)
    results['summary'] = format_result_as_list(_summary)
    return results

