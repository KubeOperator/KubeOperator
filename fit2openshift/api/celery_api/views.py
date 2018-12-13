from django.shortcuts import render


__all__ = ['celery_log_view']


def celery_log_view(request, task_id):
    return render(request, 'ansible_ui/celery_log.html', {'task_id': task_id})

